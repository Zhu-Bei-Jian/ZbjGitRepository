package gameutil

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

// CheckRedis 测试 redis 连接
func CheckRedis(addr string, password string) error {
	log := logrus.WithField("addr", addr)
	log.Info("Connect to redis...")
	conn, err := redis.Dial("tcp", addr, redis.DialPassword(password))
	if err != nil {
		log.WithError(err).Error("Connect to redis failed")
		return err
	}
	_, err = conn.Do("PING")
	conn.Close()
	if err != nil {
		log.WithError(err).Error("Connect to redis failed")
		return err
	}
	log.Info("Connect to redis succ")
	return nil
}
