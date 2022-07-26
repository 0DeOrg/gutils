package gutils

import (
	"strconv"
	"testing"
)

/**
 * @Author: lee
 * @Description:
 * @File: reflect_test
 * @Date: 2022-07-22 11:57 上午
 */

func testCall(src int, dst string) (string, int) {
	return strconv.Itoa(src) + dst, src + 100
}
func Test_reflect(t *testing.T) {
	ret := Invoke(testCall, 123, "acd")
	println(ret[0].String())
	println(ret[1].Int())

	ret = Invoke0(testCall, 123, "acd")
	println(ret[0].String())
	println(ret[1].Int())
}
