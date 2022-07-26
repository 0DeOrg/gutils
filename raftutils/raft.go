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
	"gitlab.qihangxingchen.com/qt/gutils/dumputils"
	"gitlab.qihangxingchen.com/qt/gutils/fileutils"
	"gitlab.qihangxingchen.com/qt/gutils/logutils"
	"go.uber.org/zap"
	"net"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type RaftNode struct {
	raft           *raft.Raft
	fsm            raft.FSM
	leaderNotifyCh chan bool
	logger         hclog.Logger
	bootstrap      bool
	LocalID        raft.ServerID
	LocalAddress   raft.ServerAddress
	tagLeader      int32
}

func NewRaftNode(options *options, fsm raft.FSM) (*RaftNode, error) {
	logger := hclog.Default()
	defaultCfg := raft.DefaultConfig()
	defaultCfg.LocalID = options.serverID
	defaultCfg.Logger = hclog.Default()
	notifyCh := make(chan bool, 1)
	defaultCfg.NotifyCh = notifyCh
	defaultCfg.SnapshotInterval = options.snapInterval

	tcpAddr, err := net.ResolveTCPAddr("tcp", options.bindTCPAddress)
	if nil != err {
		return nil, fmt.Errorf("NewRaftNode, an invalid tcp address: %s, err: %s", options.bindTCPAddress, err.Error())
	}

	//raft节点内部的通信通道
	transport, err := raft.NewTCPTransport(tcpAddr.String(), tcpAddr, 3, 3*time.Second, os.Stderr)
	if nil != err {
		return nil, fmt.Errorf("NewRaftNode, NewTCPTransport err: %s", err.Error())
	}

	if err = fileutils.CreateDirectoryIfNotExist(options.storeDir, os.ModePerm); nil != err {
		logutils.Fatal("NewRaftNode create dir fatal", zap.Error(err))
	}

	//快照存储，用来存储节点的快照信息
	snapshotStore, err := raft.NewFileSnapshotStoreWithLogger(options.storeDir, options.snapRetain, logger)
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

	localRaft, err := raft.NewRaft(defaultCfg, fsm, logStore, stableStore, snapshotStore, transport)
	if nil != err {
		return nil, fmt.Errorf("NewRaftNode, create raft err: %s", err.Error())
	}

	ret := &RaftNode{
		raft:           localRaft,
		fsm:            fsm,
		leaderNotifyCh: notifyCh,
		logger:         logger,
		LocalID:        defaultCfg.LocalID,
		LocalAddress:   transport.LocalAddr(),
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

	//aggressive consume notifyCh, 为避免raft堵塞需要一直消费
	go func() {
		defer dumputils.HandlePanic()
		for {
			select {
			case isLeader := <-notifyCh:
				if isLeader {
					atomic.StoreInt32(&ret.tagLeader, 1)
					logutils.Warn("node is leader", zap.String("address", string(ret.LocalAddress)))
				} else {
					atomic.StoreInt32(&ret.tagLeader, 0)
					logutils.Warn("node has lose leader", zap.String("address", string(ret.LocalAddress)))
				}
			}
		}

	}()

	return ret, nil
}

// JoinCluster
/* @Description: 加入cluster，如果当前节点不是leader 返回leader id， 即leader 监听此接口的地址（id 起了这么个作用）。
 * @param serverId string
 * @param address string
 * @return string
 * @return error
 */
func (r *RaftNode) JoinCluster(serverId string, address string) (string, error) {
	future := r.raft.AddVoter(raft.ServerID(serverId), raft.ServerAddress(address), 0, 10*time.Second)
	if err := future.Error(); nil != err {
		_, leadId := r.raft.LeaderWithID()
		return string(leadId), fmt.Errorf("JoinCluster err: %s", err.Error())
	}
	return string(r.LocalID), nil
}

// IsLeader
/* @Description: 判断当前节点是否leader节点
 * @return bool
 */
func (r *RaftNode) IsLeader() bool {
	if nil == r {
		return false
	}
	return 1 == atomic.LoadInt32(&r.tagLeader)
}

func (r *RaftNode) LeaderWithId() (string, string) {
	addr, id := r.raft.LeaderWithID()
	return string(addr), string(id)
}

func (r *RaftNode) ServerList() []raft.Server {
	future := r.raft.GetConfiguration()
	if nil != future.Error() {
		return nil
	}

	return future.Configuration().Servers
}

func (r *RaftNode) Apply(data []byte, timeout time.Duration) raft.ApplyFuture {
	return r.raft.Apply(data, timeout)
}

//func (r *RaftNode) FSMApply(data []byte) interface{} {
//	return r.fsm.Apply(&raft.Log{Data: data})
//}
