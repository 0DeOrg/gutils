package routineutils

import (
	"testing"
	"time"
)

/**
 * @Author: lee
 * @Description:
 * @File: routine_test
 * @Date: 2022-08-02 3:26 下午
 */

func test(i int, str string) {
	println(str, ": ", i)
}
func Test_routine(t *testing.T) {
	seq := NewSequenceRoutine()

	for i := 0; i < 10; i++ {
		seq.DoJob(test, i, "idx")
	}

	time.Sleep(3 * time.Second)
	for i := 0; i < 10; i++ {
		seq.DoJob(test, i, "idx2")
	}

	select {}
}
