package redislock

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
	"time"
)

const DefaultTimeout = 1 * time.Minute

var unlockScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

//分布式锁
type redisLock struct {
	resource  string
	token     string
	redisPool *redis.Pool
	timeout   time.Duration
}

func (p *redisLock) key() string {
	return fmt.Sprintf("redislock:%s", p.resource)
}

func (p *redisLock) tryLock() (ok bool, err error) {
	conn := p.redisPool.Get()
	defer conn.Close()

	status, err := redis.String(conn.Do("SET", p.key(), p.token, "EX", int64(p.timeout/time.Second), "NX"))
	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return status == "OK", nil
}

func (p *redisLock) Unlock() (err error) {
	conn := p.redisPool.Get()
	defer conn.Close()

	_, err = unlockScript.Do(conn, p.key(), p.token)
	return
}

func Lock(pool *redis.Pool, resource string, resId interface{}) (*redisLock, error) {
	var tryCount int
	for {
		lock, ok, err := tryLockWithTimeout(pool, fmt.Sprintf("%s:%v", resource, resId), DefaultTimeout)
		if err != nil {
			return nil, errors.New("error while attempting lock")
		}
		if !ok {
			tryCount++
			if tryCount > 1000 {
				return nil, errors.New("error,try too many time")
			}
			time.Sleep(time.Millisecond * 50)
			continue
		}
		return lock, nil
	}
}

func tryLockWithTimeout(pool *redis.Pool, resource string, timeout time.Duration) (lock *redisLock, ok bool, err error) {
	lock = &redisLock{
		resource:  resource,
		token:     uuid.New(),
		redisPool: pool,
		timeout:   timeout,
	}

	ok, err = lock.tryLock()

	if !ok || err != nil {
		lock = nil
	}

	return
}
