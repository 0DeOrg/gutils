package routineutils

import (
	"gitlab.qihangxingchen.com/qt/gutils"
	"log"
	"sync"
)

const max_routine_idx = 1000000000

type SequenceRoutine struct {
	cond       *sync.Cond
	awakeIdx   int
	routineIdx int
}

func NewSequenceRoutine() *SequenceRoutine {
	ret := &SequenceRoutine{
		cond: sync.NewCond(&sync.Mutex{}),
	}

	return ret
}

func (s *SequenceRoutine) DoJob(handler interface{}, params ...interface{}) {
	idx := s.routineIdx
	s.routineIdx++
	if max_routine_idx == s.routineIdx {
		s.routineIdx = 0
	}
	go func(idx int) {
		s.cond.L.Lock()
		defer s.cond.L.Unlock()
		for s.awakeIdx != idx {
			s.cond.Wait()
		}
		gutils.Invoke(handler, params...)
		s.awakeIdx = idx + 1
		idx++
		if max_routine_idx == idx {
			idx = 0
		}
		s.awakeIdx = idx
		log.Println("awakeIdx: ", idx)
		s.cond.Broadcast()
	}(idx)
}
