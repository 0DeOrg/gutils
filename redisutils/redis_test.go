package redisutils

/**
 * @Author: lee
 * @Description:
 * @File: redis_test
 * @Date: 2022-07-14 5:59 下午
 */
import (
	"github.com/0DeOrg/gutils/logutils"
	"go.uber.org/zap"
	"testing"
)

func Test_Cluster(t *testing.T) {
	cfg := ClusterCfg{
		Addrs: []string{"192.168.13.100:6379"},
		User:  "",
		Pwd:   "Q2f9YJci0dQn",
	}

	logCfg := logutils.DefaultZapConfig
	logutils.InitLogger(logCfg)
	cluster, err := NewRedisCluster(cfg)
	if nil != err {
		logutils.Fatal("NewRedisCluster err", zap.Error(err))
	}

	logutils.Info("NewRedisCluster success", zap.Any("cluster", cluster))
}
