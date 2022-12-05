package eventListener

import (
	"github.com/0DeOrg/gutils"
	"reflect"
	"sync"
	"sync/atomic"
)

/**
 * @Author: lee
 * @Description:
 * @File: event
 * @Date: 2022-10-20 3:06 下午
 */

var (
	seqMax      uint32 = 1
	listenerMap        = map[string][]*listener{}
	mtx         sync.RWMutex
)

type listener struct {
	seq     uint32
	handler interface{}
}

//除非变参， 否则handler 参数数量格式必须一致
func RegisterEvent(eventId string, handler interface{}) uint32 {
	mtx.Lock()
	defer mtx.Unlock()
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("RegisterEvent handler must be func")
	}

	ln := &listener{
		seq:     seqMax,
		handler: handler,
	}
	atomic.AddUint32(&seqMax, 1)
	list, ok := listenerMap[eventId]
	if !ok {
		list = make([]*listener, 0, 1)
	}
	listenerMap[eventId] = append(list, ln)
	return ln.seq
}

func RemoveEvent(eventId string, seq uint32) {
	mtx.Lock()
	defer mtx.Unlock()
	list, ok := listenerMap[eventId]
	if !ok {
		return
	}
	index := -1
	for i, ln := range list {
		if ln.seq == seq {
			index = i
			break
		}
	}

	if index >= 0 {
		listenerMap[eventId] = append(list[:index], list[index+1:]...)
	}
}

func TriggerEvent(eventId string, args ...interface{}) {
	mtx.RLock()
	defer mtx.RUnlock()
	list, ok := listenerMap[eventId]
	if !ok {
		return
	}

	for _, ln := range list {
		if nil != ln.handler {
			go gutils.Invoke(ln.handler, args...)
		}
	}
}
