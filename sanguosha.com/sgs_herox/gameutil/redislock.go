package gameutil

import (
	"errors"
	"github.com/everalbum/redislock"
	"github.com/garyburd/redigo/redis"
	"time"
)

//redis 分布式锁,如果没有获得，会每隔一段时间再尝试获取
func TryRedisLock(conn redis.Conn, resource string) (lock *redislock.Lock, err error) {
	for {
		lock, ok, err := redislock.TryLockWithTimeout(conn, resource, time.Second*30)
		if err != nil {
			return nil, errors.New("error while attempting lock")
		}
		if !ok {
			time.Sleep(time.Second)
			continue
		}
		return lock, nil
	}
}
