package raftutils

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2022-05-20 11:12 上午
 */

import (
	"github.com/hashicorp/raft"
	"gutils/network"
	"strconv"
	"strings"
)

//用于配置文件配置
type RaftCfg struct {
	StoreDir  string `mapstructure:"store-dir"           json:"store-dir"          yaml:"store-dir"`
	Host      string `mapstructure:"host"           json:"host"          yaml:"host"`
	Port      uint   `mapstructure:"port"           json:"port"         yaml:"port"`
	Bootstrap bool   `mapstructure:"bootstrap"           json:"bootstrap"          yaml:"bootstrap"`
}

type options struct {
	storeDir       string // store directory
	bindTCPAddress string // raft transport address
	serverID       raft.ServerID
	bootstrap      bool // start as master or not
}

func NewOptions(raftCfg *RaftCfg) *options {
	host := raftCfg.Host
	if "" == host {
		host = network.GetLocalIP()
	}

	bindTCPAddress := host + ":" + strconv.FormatUint(uint64(raftCfg.Port), 10)
	sections := strings.Split(host, ".")
	serverID := ""
	for _, v := range sections {
		serverID += "-"
		serverID += v
	}

	serverID += "-" + strconv.FormatUint(uint64(raftCfg.Port), 10)
	ret := &options{
		storeDir:       raftCfg.StoreDir,
		bindTCPAddress: bindTCPAddress,
		serverID:       raft.ServerID(serverID),
		bootstrap:      raftCfg.Bootstrap,
	}

	return ret
}
