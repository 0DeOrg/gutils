package redisutils

/**
 * @Author: lee
 * @Description:
 * @File: redis-cluster
 * @Date: 2022/3/10 2:38 下午
 */

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type ClusterCfg struct {
	Addrs []string `json:"addrs"     yaml:"addrs"   mapstructure:"addrs"`
	User  string   `json:"user"     yaml:"user"   mapstructure:"user"`
	Pwd   string   `json:"pwd"     yaml:"pwd"   mapstructure:"pwd"`
}

type RedisCluster struct {
	cfg    ClusterCfg
	Client *redis.ClusterClient
}

func NewRedisCluster(cfg ClusterCfg) (*RedisCluster, error) {

	opt := redis.ClusterOptions{
		Addrs:    cfg.Addrs,
		Username: cfg.User,
		Password: cfg.Pwd,
	}

	client := redis.NewClusterClient(&opt)
	_, err := client.Ping(context.TODO()).Result()
	if nil != err {
		return nil, fmt.Errorf("NewRedisCluster fatal, ping failed, err: %s", err.Error())
	}

	ret := &RedisCluster{
		cfg:    cfg,
		Client: client,
	}

	return ret, nil
}

func (r *RedisCluster) ReleaseResource() (err error) {
	if nil != r.Client {
		err = r.Client.Close()
		return
	}

	return
}
