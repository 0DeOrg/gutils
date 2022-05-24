package raftutils

/**
 * @Author: lee
 * @Description:
 * @File: raft
 * @Date: 2022-05-20 11:12 上午
 */

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"net"
	"os"
	"path/filepath"
	"time"
)

type RaftNode struct {
	raft           *raft.Raft
	fsm            raft.FSM
	leaderNotifyCh chan bool
	logger         hclog.Logger
	bootstrap      bool
}

func NewRaftNode(options *options, applyCallback FSMApplyFunc) (*RaftNode, error) {
	logger := hclog.Default()
	defaultCfg := raft.DefaultConfig()
	defaultCfg.LocalID = options.serverID
	defaultCfg.Logger = hclog.Default()
	notifyCh := make(chan bool, 1)
	defaultCfg.NotifyCh = notifyCh

	tcpAddr, err := net.ResolveTCPAddr("tcp", options.bindTCPAddress)
	if nil != err {
		return nil, fmt.Errorf("NewRaftNode, an invalid tcp address: %s, err: %s", options.bindTCPAddress, err.Error())
	}

	//raft节点内部的通信通道
	transport, err := raft.NewTCPTransport(tcpAddr.String(), tcpAddr, 3, 3*time.Second, os.Stderr)
	if nil != err {
		return nil, fmt.Errorf("NewRaftNode, NewTCPTransport err: %s", err.Error())
	}

	//快照存储，用来存储节点的快照信息
	snapshotStore, err := raft.NewFileSnapshotStoreWithLogger(options.storeDir, 2, logger)
	if nil != err {
		return nil, fmt.Errorf("NewRaftNode, snapshot store err: %s", err.Error())
	}

	//用来存储raft的日志
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(options.storeDir, "raft-log.bolt"))
	if err != nil {
		return nil, fmt.Errorf("NewRaftNode, logStore err: %s", err.Error())
	}

	//稳定存储，用来存储raft集群的节点信息等
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(options.storeDir, "raft-stable.bolt"))
	if err != nil {
		return nil, fmt.Errorf("NewRaftNode, stableStore err: %s", err.Error())
	}

	fsm := NewFSM(applyCallback)

	localRaft, err := raft.NewRaft(defaultCfg, fsm, logStore, stableStore, snapshotStore, transport)
	if nil != err {
		return nil, fmt.Errorf("NewRaftNode, create raft err: %s", err.Error())
	}

	ret := &RaftNode{
		raft:           localRaft,
		fsm:            fsm,
		leaderNotifyCh: notifyCh,
		logger:         logger,
	}

	//是否引导启动，只有一个是作为引导启动的
	if options.bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      defaultCfg.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}

		localRaft.BootstrapCluster(configuration)
		ret.bootstrap = true
	}

	return ret, nil
}

func (r *RaftNode) JoinCluster(serverId string, address string) error {
	future := r.raft.AddVoter(raft.ServerID(serverId), raft.ServerAddress(address), 0, 10*time.Second)
	if err := future.Error(); nil != err {
		return fmt.Errorf("JoinCluster err: %s", err.Error())
	}
	return nil
}
