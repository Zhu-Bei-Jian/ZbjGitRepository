package userdb

import (
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"strconv"
	"time"

	"sanguosha.com/sgs_herox/gameutil/cache"
)

const (
	autoExpireDuration      = time.Hour + 5*time.Minute
	autoSyncDuration        = time.Minute * 5
	autoSyncDurationDevelop = time.Second * 5
)

type shard struct {
	ch     chan func()
	cache  *cache.Cache
	config *conf.GameConfig
}

func newShard(db *DB, chanLen int, dev bool) *shard {
	s := new(shard)
	s.ch = make(chan func(), chanLen)

	load := func(key string) (interface{}, error) {
		userid, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return nil, err
		}
		return db.loadUser(userid)
	}
	sync := func(key string, obj interface{}) error {
		u := obj.(*User)
		return u.sync()
	}

	worker := func(f func()) {
		db.rw.RLock()
		if !db.closed {
			s.ch <- f
		}
		//closed := db.closed
		db.rw.RUnlock()
		//if closed {
		//	return
		//}
		//s.ch <- f
	}
	s.cache = cache.NewCache(load, sync, worker)

	var autoSyncDuration time.Duration = autoSyncDuration
	if dev {
		autoSyncDuration = autoSyncDurationDevelop
	}
	s.cache.AutoSync(autoSyncDuration)
	s.cache.AutoExpire(autoExpireDuration)

	return s
}

func (s *shard) GetCacheInfo() int64 {
	return s.cache.GetCacheInfo()
}

func (s *shard) GetAllUser() (all []*User) {
	objs := s.cache.GetAll()
	for _, obj := range objs {
		all = append(all, obj.(*User))
	}
	return all
}

func (s *shard) checkUser(userid uint64) *User {
	key := strconv.FormatUint(userid, 10)
	obj := s.cache.Check(key)
	if obj == nil {
		return nil
	}
	return obj.(*User)
}

func (s *shard) getUser(userid uint64) (*User, error) {
	key := strconv.FormatUint(userid, 10)
	obj, err := s.cache.Get(key, false)
	if err != nil {
		return nil, err
	}
	return obj.(*User), nil
}

func (s *shard) recoverUser(userid uint64) (*User, error) {
	key := strconv.FormatUint(userid, 10)
	obj, err := s.cache.Get(key, true)
	if err != nil {
		return nil, err
	}
	return obj.(*User), nil
}

func (s *shard) cache2DB(userid uint64) {
	key := strconv.FormatUint(userid, 10)
	s.cache.Cache2DB(key)
}

func (s *shard) sync() error {
	return s.cache.Sync()
}
