package structutils

import "testing"

/**
 * @Author: lee
 * @Description:
 * @File: queue_test
 * @Date: 2022-08-11 11:46 上午
 */
func Test_QueueFIFO(t *testing.T) {
	q := NewQueueFIFO(5)
	for i := 0; i < 20; i++ {
		q.Push(i)
		q.Print()
		q.PrintReverse()
	}
}
