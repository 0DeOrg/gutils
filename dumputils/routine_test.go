package dumputils

/**
 * @Author: lee
 * @Description:
 * @File: routine_test
 * @Date: 2022/2/21 10:56 上午
 */
import "testing"

func Test_Dump(t *testing.T) {
	defer HandlePanic()

	panic("该结束了")

}
