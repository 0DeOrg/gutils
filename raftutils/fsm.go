package raftutils

/**
 * @Author: lee
 * @Description:
 * @File: fsm
 * @Date: 2022-05-20 4:11 下午
 */

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"io"
)

type FSM struct {
	cache     *logCache
	applyCall FSMApplyFunc
}

type FSMApplyFunc func() interface{}

func NewFSM(applyCall FSMApplyFunc) *FSM {
	ret := &FSM{
		cache:     newLogCache(),
		applyCall: applyCall,
	}

	return ret
}

type logEntryData struct {
	Key   string
	Value string
}

var _ raft.FSM = (*FSM)(nil)

// Apply is called once a log entry is committed by a majority of the cluster.
//
// Apply should apply the log to the FSM. Apply must be deterministic and
// produce the same result on all peers in the cluster.
//
// The returned value is returned to the client as the ApplyFuture.Response.
func (fsm *FSM) Apply(log *raft.Log) interface{} {
	logEntry := logEntryData{}
	_ = json.Unmarshal(log.Data, &logEntry)

	fsm.cache.Set(logEntry.Key, logEntry.Value)
	if nil != fsm.applyCall {
		return fsm.applyCall()
	}

	return nil
}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &FSMSnapshot{cache: fsm.cache}, nil
}

func (fsm *FSM) Restore(snapshot io.ReadCloser) error {
	return fsm.cache.UnMarshal(snapshot)
}

type FSMSnapshot struct {
	cache *logCache
}

var _ raft.FSMSnapshot = (*FSMSnapshot)(nil)

// Persist should dump all necessary state to the WriteCloser 'sink',
// and call sink.Close() when finished or call sink.Cancel() on error.
func (s *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	data, err := s.cache.Marshal()
	if nil != err {
		return fmt.Errorf("FSMSnapshot Persist Marshal err: %s", err.Error())
	}

	if _, err = sink.Write(data); err != nil {
		_ = sink.Cancel()
		return fmt.Errorf("FSMSnapshot Persist Write err: %s", err.Error())
	}

	if err = sink.Close(); err != nil {
		_ = sink.Cancel()
		return fmt.Errorf("FSMSnapshot Persist Marshal err: %s", err.Error())
	}

	return nil
}

// Release
/* @Description: 快照处理完成后的回调
 */
func (s *FSMSnapshot) Release() {

}
