package redisutils
/**
 * @Author: lee
 * @Description:
 * @File: redis
 * @Date: 2021/10/11 7:54 下午
 */

import (
	"github.com/garyburd/redigo/redis"
	"github.com/juju/ratelimit"
	"time"
)

type RedisConfig struct {
	DB			int 		`json:"db"     yaml:"db"   mapstructure:"db"`
	Address 	string 		`json:"address"     yaml:"address"   mapstructure:"address"`
	Pwd 		string 		`json:"pwd"     yaml:"pwd"   mapstructure:"pwd"`
	MaxIdle 	int 		`json:"maxIdle"     yaml:"max-idle"   mapstructure:"max-idle"`
	IdleTimeout int 		`json:"idleTimeout"     yaml:"idle-timeout"   mapstructure:"idle-timeout"`
}

type RedisAgent struct {
	config RedisConfig
	pool *redis.Pool
	bucket *ratelimit.Bucket
}

func NewRedisAgent(config RedisConfig) *RedisAgent {
	pool := redis.Pool{
		MaxIdle: config.MaxIdle,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error){
			c, err := redis.Dial("tcp", config.Address)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", config.Pwd); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", config.DB); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}

	ret := RedisAgent{
		config: config,
		pool: &pool,
		bucket: ratelimit.NewBucketWithQuantum(100 * time.Millisecond, 50, 1),
	}

	return &ret
}

func (r *RedisAgent)Ping() error{
	c := r.GetConn()
	if nil != c.Err() {
		_, err := c.Do("PING")
		if nil != err {
			return err
		}

		return nil
	}

	return c.Err()
}

func (r *RedisAgent)GetConn() redis.Conn {
	return r.pool.Get()
}
