package cache

import (
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

// LoadFunc 载入数据委托函数
type LoadFunc func(key string) (interface{}, error)

// SyncFunc 数据同步委托函数
type SyncFunc func(key string, obj interface{}) error

type cacheObject struct {
	obj         interface{}
	expireTimer *time.Timer
	syncTimer   *time.Timer
}

// Cache 处理缓存逻辑.
type Cache struct {
	load               LoadFunc
	sync               SyncFunc
	worker             func(func())
	ableAutoSync       bool
	autoSyncDuration   time.Duration
	ableAutoExpire     bool
	autoExpireDuration time.Duration

	objs map[string]*cacheObject

	count int64
}

// NewCache 创建一个 Cache
func NewCache(load LoadFunc, sync SyncFunc, worker func(func())) *Cache {
	c := new(Cache)
	c.load = load
	c.sync = sync
	c.worker = worker
	c.objs = make(map[string]*cacheObject)
	return c
}

// AutoSync 开启自动同步, 并设置同步时间间隔
func (c *Cache) AutoSync(d time.Duration) {
	c.ableAutoSync = true
	c.autoSyncDuration = d
}

// AutoExpire 开启自动过期处理, 并设置自动过期时间
// TODO 目前策略比较简单, 每次 Get 后重置一个倒计时, 到期移除缓存对象. 后期考虑实现 LRU 策略
func (c *Cache) AutoExpire(d time.Duration) {
	c.ableAutoExpire = true
	c.autoExpireDuration = d
}

func (c *Cache) GetCacheInfo() int64 {
	v := atomic.LoadInt64(&c.count)
	return v
}

//All 获取内存存在的所有对象
func (c *Cache) GetAll() (all []interface{}) {
	count := len(c.objs)
	if count == 0 {
		return
	}
	all = make([]interface{}, count)
	i := 0
	for _, v := range c.objs {
		all[i] = v.obj
		i++
	}
	return all
}

//Check 获取内存存在的对象
func (c *Cache) Check(key string) interface{} {
	if v, ok := c.objs[key]; ok {
		return v.obj
	}
	return nil
}

// Get 获取一个 key 对应的对象, 不存在时载入, 载入后缓存
func (c *Cache) Get(key string, recover bool) (interface{}, error) {
	var cacheObj *cacheObject
	if v, ok := c.objs[key]; ok {
		cacheObj = v

		if recover {
			if cacheObj.syncTimer != nil {
				cacheObj.syncTimer.Stop()
			}
			if cacheObj.expireTimer != nil {
				cacheObj.expireTimer.Stop()
			}
			atomic.AddInt64(&c.count, -1)
			delete(c.objs, key)
			logrus.Debug("Cache Recover delete ", key)
			cacheObj = nil
		}
	}
	if cacheObj == nil {
		obj, err := c.load(key)
		if err != nil {
			return nil, err
		}
		cacheObj = &cacheObject{
			obj: obj,
		}
		c.objs[key] = cacheObj
		atomic.AddInt64(&c.count, 1)
		logrus.Debug("Cache Get load ", key)
	}

	// 按固定时间间隔同步
	if c.ableAutoSync {
		if cacheObj.syncTimer == nil {
			var f func()
			f = func() {
				c.worker(func() {
					if c.objs[key] != cacheObj {
						return
					}
					cacheObj.syncTimer = time.AfterFunc(c.autoSyncDuration, f)
					c.sync(key, cacheObj.obj)
				})
			}
			cacheObj.syncTimer = time.AfterFunc(c.autoSyncDuration, f)
		}
	}
	// 过期删除
	if c.ableAutoExpire {
		if cacheObj.expireTimer == nil {
			cacheObj.expireTimer = time.AfterFunc(c.autoExpireDuration, func() {
				c.worker(func() {
					if c.objs[key] != cacheObj {
						return
					}
					if cacheObj.syncTimer != nil {
						cacheObj.syncTimer.Stop()
					}
					c.sync(key, cacheObj.obj)
					atomic.AddInt64(&c.count, -1)
					delete(c.objs, key)
					logrus.Debug("Cache Expire delete ", key)
				})
			})
		} else {
			// 每一次 Get 重置过期倒计时
			cacheObj.expireTimer.Reset(c.autoExpireDuration)
		}
	}

	return cacheObj.obj, nil
}

// Find 查找已经缓存的对象.
func (c *Cache) Find(key string) (interface{}, bool) {
	if v, ok := c.objs[key]; ok {
		return v.obj, true
	}
	return nil, false
}

//使对象立即储存到数据库,并从内存中删除
func (c *Cache) Cache2DB(key string) bool {
	if v, ok := c.objs[key]; ok {
		if v.syncTimer != nil {
			v.syncTimer.Stop()
		}
		if v.expireTimer != nil {
			v.expireTimer.Stop()
		}
		c.sync(key, v.obj)
		atomic.AddInt64(&c.count, -1)
		delete(c.objs, key)
		logrus.Debug("Cache Cache2DB delete ", key)
		return true
	}
	return false
}

// Sync 同步所有数据
func (c *Cache) Sync() error {
	var err error
	for key, cobj := range c.objs {
		e := c.sync(key, cobj.obj)
		if e != nil {
			err = e
		}
	}
	return err
}
