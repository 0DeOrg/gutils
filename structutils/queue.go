package structutils

/**
 * @Author: lee
 * @Description: 固定长度的队列，如果队列满了则将头部的删除
 * @File: queue
 * @Date: 2022-08-11 10:49 上午
 */
import (
	"errors"
	"fmt"
	"sync"
)

/* @Description: 固定长度的队列，如果队列满了则将z
 */

type QueueFIFO struct {
	Capacity int           `json:"c"`
	H        int           `json:"h"`
	T        int           `json:"t"`
	Length   int           `json:"l"`
	Cache    []interface{} `json:"C"` //最好是相同类型的
	mtx      sync.RWMutex
}

func NewQueueFIFO(capacity int) *QueueFIFO {
	ret := &QueueFIFO{
		Capacity: capacity,
		Cache:    make([]interface{}, capacity, capacity),
	}

	return ret
}

// Push
/* @Description: 队列尾插入数据
 * @param ele interface{}
 */
func (q *QueueFIFO) Push(ele interface{}) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.Length > 0 {
		q.T = (q.T + 1) % q.Capacity
		if q.T == q.H {
			q.H = (q.H + 1) % q.Capacity
		}
	}
	q.Cache[q.T] = ele

	q.Length++
	if q.Length > q.Capacity {
		q.Length = q.Capacity
	}
}

func (q *QueueFIFO) UpdateFromTail(idx int, value interface{}) bool {
	if idx >= q.Length {
		return false
	}

	realIdx := (q.T - idx + q.Capacity) % q.Capacity
	q.Cache[realIdx] = value

	return true
}

func (q *QueueFIFO) UpdateFromHead(idx int, value interface{}) bool {
	if idx >= q.Length {
		return false
	}

	realIdx := (q.H + idx) % q.Capacity
	q.Cache[realIdx] = value

	return true
}

// Range
/* @Description: 从队列尾开始遍历
 * @param fn func(int, interface{}) bool
 */
func (q *QueueFIFO) Range(fn func(int, interface{}) bool) {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	for i := 0; i < q.Length; i++ {
		idx := (q.T - i + q.Capacity) % q.Capacity
		if !fn(i, q.Cache[idx]) {
			break
		}
	}
}

// ReverseRange
/* @Description: 从队列头开始遍历
 * @param fn func(int, interface{}) bool
 */
func (q *QueueFIFO) ReverseRange(fn func(int, interface{}) bool) {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	for i := 0; i < q.Length; i++ {
		idx := (q.H + i) % q.Capacity
		if !fn(i, q.Cache[idx]) {
			break
		}
	}
}

func (q *QueueFIFO) Resize(capacity int) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	if capacity == q.Capacity {
		return
	}
	cache := make([]interface{}, capacity, capacity)
	offset := q.Length - capacity
	tail := 0
	foreach := func(idx int, value interface{}) bool {
		if idx < offset {
			return true
		}
		cache[tail] = value
		tail++

		return true
	}
	for i := 0; i < q.Length; i++ {
		idx := (q.H + i) % q.Capacity
		if !foreach(i, q.Cache[idx]) {
			break
		}
	}
	q.H = 0
	q.T = tail - 1
	q.Capacity = capacity
	//更新长度
	if q.Length > capacity {
		q.Length = capacity
	}

	q.Cache = cache
}

func (q *QueueFIFO) Clear() {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.Length = 0
	q.H = 0
	q.T = 0
}

func (q *QueueFIFO) IsFull() bool {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	return q.Length == q.Capacity
}

func (q *QueueFIFO) IsEmpty() bool {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	return 0 == q.Length
}

func (q *QueueFIFO) Len() int {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	return q.Length
}

func (q *QueueFIFO) Head() (interface{}, error) {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	if 0 == q.Length {
		return nil, errors.New("queue is empty")
	}
	return q.Cache[q.H], nil
}

func (q *QueueFIFO) Tail() (interface{}, error) {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	if 0 == q.Length {
		return nil, errors.New("queue is empty")
	}
	return q.Cache[q.T], nil
}

func (q *QueueFIFO) Print() {
	fmt.Print("print  : ")
	q.Range(func(idx int, value interface{}) bool {
		fmt.Print(value, ", ")
		return true
	})
	fmt.Println()
}
func (q *QueueFIFO) PrintReverse() {
	fmt.Print("reverse print: ")
	q.ReverseRange(func(idx int, value interface{}) bool {
		fmt.Print(value, ", ")
		return true
	})
	fmt.Println()
}
