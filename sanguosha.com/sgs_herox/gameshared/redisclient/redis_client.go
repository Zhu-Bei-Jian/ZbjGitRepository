package redisclient

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/gameshared/config"
	"time"
)

type RedisClient struct {
	pool *redis.Pool
}

func (p *RedisClient) GetPool() *redis.Pool {
	return p.pool
}

func (p *RedisClient) Init(redisCfg *config.RedisConfig) error {
	p.pool = &redis.Pool{
		MaxIdle:     redisCfg.MaxIdle,
		MaxActive:   redisCfg.Max,
		IdleTimeout: 300 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			logrus.WithFields(logrus.Fields{
				"address": redisCfg.Addr,
			}).Info("正在连接 redis")
			c, err := redis.Dial("tcp", redisCfg.Addr)
			if err != nil {
				return nil, err
			}
			if redisCfg.Password != "" {
				if _, err := c.Do("AUTH", redisCfg.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			logrus.Info("连接 redis 成功")
			return c, nil
		},
	}

	// 测试redis连接
	{
		conn := p.pool.Get()
		if err := conn.Err(); err != nil {
			return err
		}
		_, err := conn.Do("PING")
		conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *RedisClient) Close() error {
	return p.pool.Close()
}

func (p *RedisClient) GetConn() redis.Conn {
	return p.pool.Get()
}

func (p *RedisClient) Exec(cmd string, keyAndArgs ...interface{}) (interface{}, error) {
	conn := p.pool.Get()
	// put connection back to redis pool
	defer conn.Close()

	reply, err := conn.Do(cmd, redis.Args{}.Add().AddFlat(keyAndArgs)...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"command":    cmd,
			"keyAndArgs": keyAndArgs,
		}).WithError(err).Error("redis.Exec error")
		return nil, err
	}
	return reply, nil
}

func (p *RedisClient) Set(key string, value string) (string, error) {
	return redis.String(p.Exec("SET", key, value))
}

func (p *RedisClient) SetEx(key string, value interface{}, expire time.Duration) (string, error) {
	return redis.String(p.Exec("SET", key, value, "EX", int64(expire.Seconds())))
}

func (p *RedisClient) SetNx(key string, value interface{}, expire time.Duration) (string, error) {
	return redis.String(p.Exec("SET", key, value, "EX", int64(expire.Seconds()), "NX"))
}

func (p *RedisClient) Get(key string) (string, error) {
	return redis.String(p.Exec("GET", key))
}

func (p *RedisClient) GetBytes(key string) ([]byte, error) {
	return redis.Bytes(p.Exec("GET", key))
}

func (p *RedisClient) MGET(key ...string) ([]string, error) {
	args := redis.Args{}.AddFlat(key)
	return redis.Strings(p.Exec("MGET", args...))
}

func (p *RedisClient) MGetBytes(key ...string) ([][]byte, error) {
	args := redis.Args{}.AddFlat(key)
	return redis.ByteSlices(p.Exec("MGET", args...))
}

func (p *RedisClient) Del(key string) (int64, error) {
	return redis.Int64(p.Exec("DEL", key))
}

func (p *RedisClient) Keys(pattern string) ([]string, error) {
	return redis.Strings(p.Exec("KEYS", pattern))
}

func (p *RedisClient) Incr(key string) (int64, error) {
	return redis.Int64(p.Exec("INCR", key))
}

func (p *RedisClient) Decr(key string) (int64, error) {
	return redis.Int64(p.Exec("DECR", key))
}

func (p *RedisClient) IncrBy(key string, num int64) (int64, error) {
	return redis.Int64(p.Exec("INCRBY", key, num))
}

func (p *RedisClient) HGet(key string, field interface{}) (string, error) {
	return redis.String(p.Exec("HGET", key, field))
}

func (p *RedisClient) HGetBytes(key string, field interface{}) ([]byte, error) {
	return redis.Bytes(p.Exec("HGET", key, field))
}

func (p *RedisClient) HGetInt64(key string, field interface{}) (int64, error) {
	return redis.Int64(p.Exec("HGET", key, field))
}

func (p *RedisClient) HSet(key, field interface{}, value interface{}) (int64, error) {
	reply, err := p.Exec("HSET", key, field, value)
	return redis.Int64(reply, err)
}

func (p *RedisClient) HKeys(key string) ([]string, error) {
	return redis.Strings(p.Exec("HKEYS", key))
}

func (p *RedisClient) HSETNX(key, field interface{}, value string) (int64, error) {
	return redis.Int64(p.Exec("HSETNX", key, field, value))
}

func (p *RedisClient) HDel(key string, field interface{}) (int64, error) {
	return redis.Int64(p.Exec("HDEL", key, field))
}

func (p *RedisClient) HDels(key string, field ...string) ([]string, error) {
	args := redis.Args{}.Add(key).AddFlat(field)
	return redis.Strings(p.Exec("HDEL", args...))
}

func (p *RedisClient) HExists(key string, field string) (bool, error) {
	i64, err := redis.Int64(p.Exec("HEXISTS", key, field))
	if err != nil {
		return false, err
	}
	return i64 == 1, nil
}

func (p *RedisClient) HLen(key string) (int64, error) {
	return redis.Int64(p.Exec("HLEN", key))
}

func (p *RedisClient) LLen(key string) (int64, error) {
	return redis.Int64(p.Exec("LLEN", key))
}

func (p *RedisClient) HIncrBy(key string, field interface{}, increment int64) (int64, error) {
	return redis.Int64(p.Exec("HINCRBY", key, field, increment))
}

func (p *RedisClient) HMSet(key string, pairs ...string) (string, error) {
	args := redis.Args{}.Add(key).AddFlat(pairs)
	return redis.String(p.Exec("HMSET", args...))
}

func (p *RedisClient) HMGet(key string, field ...string) ([]string, error) {
	args := redis.Args{}.Add(key).AddFlat(field)
	return redis.Strings(p.Exec("HMGET", args...))
}

func (p *RedisClient) HGetAll(key string) (map[string]string, error) {

	return redis.StringMap(p.Exec("HGETALL", key))
}

func (p *RedisClient) ZCard(key string) (int, error) {
	return redis.Int(p.Exec("ZCARD", key))
}

func (p *RedisClient) ZAdd(key string, score int64, member interface{}) (int64, error) {
	return redis.Int64(p.Exec("ZADD", key, score, member))
}

func (p *RedisClient) Expire(key string, expire int64) (uint64, error) {
	return redis.Uint64(p.Exec("EXPIRE", key, expire))
}

func (p *RedisClient) ZAddFloat(key string, score float64, member string) (int64, error) {
	return redis.Int64(p.Exec("ZADD", key, score, member))
}

//要测试可用性
func (p *RedisClient) ZAdds(key string, pairs ...interface{}) (int64, error) {
	keyAndArgs := make([]interface{}, 0, 1+len(pairs))
	keyAndArgs = append(keyAndArgs, key)
	keyAndArgs = append(keyAndArgs, pairs...)
	return redis.Int64(p.Exec("ZADD", keyAndArgs...))
}

func (p *RedisClient) ZRem(key string, member interface{}) (int64, error) {
	return redis.Int64(p.Exec("ZREM", key, member))
}

func (p *RedisClient) ZIncrBy(key string, increment int32, member string) (string, error) {
	return redis.String(p.Exec("ZINCRBY", key, increment, member))
}

func (p *RedisClient) ZIncrByFloat(key string, increment float64, member string) (string, error) {
	return redis.String(p.Exec("ZINCRBY", key, increment, member))
}

//获取member成员的score值
func (p *RedisClient) ZScore(key string, member string) (int64, error) {
	return redis.Int64(p.Exec("ZSCORE", key, member))
}

func (p *RedisClient) ZScoreFloat(key string, member string) (float64, error) {
	return redis.Float64(p.Exec("ZSCORE", key, member))
}

//从小到大排序的名次
func (p *RedisClient) ZRank(key string, member string) (int64, error) {
	rank, err := redis.Int64(p.Exec("ZRANK", key, member))
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}

//从大到小排序的名次
func (p *RedisClient) ZRevRank(key string, member string) (int64, error) {
	rank, err := redis.Int64(p.Exec("ZREVRANK", key, member))
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}

func (p *RedisClient) ZRange(key string, start, stop int32) ([]string, error) {
	return redis.Strings(p.Exec("ZRANGE", key, start, stop))
}

func (p *RedisClient) ZRangeWithScores(key string, start, stop int32) ([]string, error) {
	return redis.Strings(p.Exec("ZRANGE", key, start, stop, "WITHSCORES"))
}

func (p *RedisClient) ZRevRange(key string, start, stop int32) ([]string, error) {
	return redis.Strings(p.Exec("ZREVRANGE", key, start, stop))
}

func (p *RedisClient) ZRevRangeWithScores(key string, start, stop int32) ([]string, error) {
	return redis.Strings(p.Exec("ZREVRANGE", key, start, stop, "WITHSCORES"))
}

func (p *RedisClient) ZRangeByScoreWithScores(key string, min, max int64) ([]string, error) {
	return redis.Strings(p.Exec("ZRANGEBYSCORE", key, min, max, "WITHSCORES"))
}

//未验证可用性
func (p *RedisClient) ZRangeByScore(key string, min, max int64) ([]string, error) {
	conn := p.pool.Get()
	defer conn.Close()

	vs, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max))
	if err != nil {
		return nil, err
	}
	nLen := len(vs)
	if nLen == 0 {
		return nil, nil
	}

	ivs := make([]string, nLen)
	err = redis.ScanSlice(vs, &ivs)
	if err != nil {
		return nil, err
	}
	return ivs, nil
}

//未验证可用性
func (p *RedisClient) ZRemRangeByScore(key string, min, max int64) error {
	conn := p.pool.Get()
	defer conn.Close()

	_, err := redis.Values(conn.Do("ZREMRANGEBYSCORE", key, min, max))
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisClient) ZRemRangeByRank(key string, min, max int64) error {
	conn := p.pool.Get()
	defer conn.Close()

	_, err := redis.Values(conn.Do("ZREMRANGEBYRANK", key, min, max))
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisClient) RPush(key string, value string) (int64, error) {
	return redis.Int64(p.Exec("RPUSH", key, value))
}

func (p *RedisClient) LPop(key string) (string, error) {
	reply, err := p.Exec("LPOP", key)
	return redis.String(reply, err)
}

func (p *RedisClient) SAdd(key string, members ...string) (int64, error) {
	return redis.Int64(p.Exec("SADD", key, members))
}

func (p *RedisClient) SMembers(key string) ([]string, error) {
	return redis.Strings(p.Exec("SMEMBERS", key))
}

func (p *RedisClient) SRem(key string, member string) (int64, error) {
	return redis.Int64(p.Exec("SREM", key, member))
}

////push推送事件
//func (p *RedisClient) PushGetUIEvent(pushTime int64, userID uint64, cid string, eventType gamedef.GetUIEvent_EventType) {
//	event := &gamedef.GetUIEvent{
//		PushTime:  pushTime,
//		UserID:    userID,
//		Cid:       cid,
//		EventType: eventType,
//	}
//}
//
//func (p *RedisClient) GetLastestGetUIEvent() *gamedef.GetUIEvent {
//
//}
//
//func (p *RedisClient) PopLatestGetUIEvent() *gamedef.GetUIEvent {
//
//}
