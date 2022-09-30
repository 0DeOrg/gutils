package structutils

import (
	"errors"
	"fmt"
	"sync"
)

/**
 * @Author: lee
 * @Description:
 * @File: queue
 * @Date: 2022-08-11 10:49 上午
 */

/* @Description: 固定长度的队列，如果队列满了则将z
 */

type QueueFIFO struct {
	capacity int
	head     int
	tail     int
	len      int
	cache    []interface{} //最好是相同类型的
	mtx      sync.RWMutex
}

func NewQueueFIFO(capacity int) *QueueFIFO {
	ret := &QueueFIFO{
		capacity: capacity,
		cache:    make([]interface{}, capacity, capacity),
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

	q.cache[q.tail] = ele
	q.tail = (q.tail + 1) % q.capacity
	if q.len == q.capacity {
		q.head = (q.head + 1) % q.capacity
	}

	q.len++
	if q.len > q.capacity {
		q.len = q.capacity
	}
}

// Range
/* @Description: 从队列尾开始遍历
 * @param fn func(int, interface{}) bool
 */
func (q *QueueFIFO) Range(fn func(int, interface{}) bool) {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	for i := 0; i < q.len; i++ {
		idx := (q.tail - 1 - i + q.capacity) % q.capacity
		if !fn(i, q.cache[idx]) {
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
	for i := 0; i < q.len; i++ {
		idx := (q.head + i) % q.capacity
		if !fn(i, q.cache[idx]) {
			break
		}
	}
}

func (q *QueueFIFO) Clear() {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.len = 0
	q.head = 0
	q.tail = 0
}

func (q *QueueFIFO) IsFull() bool {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	return q.len == q.capacity
}

func (q *QueueFIFO) IsEmpty() bool {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	return 0 == q.len
}

func (q *QueueFIFO) Len() int {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	return q.len
}

func (q *QueueFIFO) Head() (interface{}, error) {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	if 0 == q.len {
		return nil, errors.New("queue is empty")
	}
	return q.cache[q.head], nil
}

func (q *QueueFIFO) Tail() (interface{}, error) {
	q.mtx.RLock()
	defer q.mtx.RUnlock()
	if 0 == q.len {
		return nil, errors.New("queue is empty")
	}
	return q.cache[(q.tail-1+q.capacity)%q.capacity], nil
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
	fmt.Print("reverse: ")
	q.ReverseRange(func(idx int, value interface{}) bool {
		fmt.Print(value, ", ")
		return true
	})
	fmt.Println()
}
