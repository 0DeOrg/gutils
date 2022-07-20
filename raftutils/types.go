package raftutils

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2022-05-20 11:12 上午
 */

import (
	"fmt"
	"github.com/hashicorp/raft"
	"gutils/network"
	"path"
	"strconv"
)

//用于配置文件配置
type RaftCfg struct {
	StoreDir  string `mapstructure:"store-dir"           json:"store-dir"          yaml:"store-dir"`
	Host      string `mapstructure:"host"           json:"host"          yaml:"host"`
	Port      uint   `mapstructure:"port"           json:"port"         yaml:"port"`
	Bootstrap bool   `mapstructure:"bootstrap"           json:"bootstrap"          yaml:"bootstrap"`
	HttpPort  uint   `mapstructure:"http-port"           json:"http-port"         yaml:"http-port"`
}

type options struct {
	storeDir       string        // store directory
	bindTCPAddress string        // raft transport address
	serverID       raft.ServerID // serverID 直接用监听加入集群地址作为id，方便leader轮换后找到新的leader监听地址
	bootstrap      bool          // start as master or not
}

func (o *options) ServerID() raft.ServerID {
	return o.serverID
}

func NewOptions(raftCfg *RaftCfg, httpPort uint) (*options, error) {
	host := raftCfg.Host
	if "" == host {
		host = network.GetLocalIP()
	}

	if 0 == raftCfg.HttpPort && 0 == httpPort {
		return nil, fmt.Errorf("both raft cfg http port and system port are zero")
	}

	serverId := host
	if 0 == raftCfg.HttpPort {
		serverId += ":" + strconv.FormatUint(uint64(httpPort), 10)
	} else {
		serverId += ":" + strconv.FormatUint(uint64(raftCfg.HttpPort), 10)
	}

	bindTCPAddress := host + ":" + strconv.FormatUint(uint64(raftCfg.Port), 10)

	ret := &options{
		storeDir:       path.Join(raftCfg.StoreDir, serverId),
		bindTCPAddress: bindTCPAddress,
		serverID:       raft.ServerID(serverId),
		bootstrap:      raftCfg.Bootstrap,
	}

	return ret, nil
}
