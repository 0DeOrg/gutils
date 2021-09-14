package convert

import "github.com/spf13/cast"

/**
 * @Author: lee
 * @Description:
 * @File: cast_e
 * @Date: 2021/9/2 6:11 下午
 */


func ToStringE(i interface{}) (string, error) {
	return cast.ToStringE(i)
}

func ToIntE(i interface{}) (int, error) {
	return cast.ToIntE(i)
}

func ToUintE(i interface{}) (uint, error) {
	return cast.ToUintE(i)
}

func ToStringMapStringE(i interface{}) (map[string]string, error) {
	return cast.ToStringMapStringE(i)
}


