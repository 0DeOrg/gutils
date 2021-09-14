package judge
/**
 * @Author: lee
 * @Description:
 * @File: judge
 * @Date: 2021/9/14 11:58 上午
 */


import (
	"reflect"
)

func IsPtr(i interface{}) bool {
	refType := reflect.TypeOf(i)
	return reflect.Ptr == refType.Kind()
}

func IsStructPtr(i interface{}) bool {
	refType := reflect.TypeOf(i)
	if reflect.Ptr != refType.Kind() {
		return false
	} else {
		refType = refType.Elem()
	}

	return reflect.Struct == refType.Kind()
}