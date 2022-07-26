package gutils

import "reflect"

/**
 * @Author: lee
 * @Description:
 * @File: reflect
 * @Date: 2022-07-22 11:51 上午
 */

func Invoke(fn interface{}, params ...interface{}) []reflect.Value {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return nil
	}

	fv := reflect.ValueOf(fn)

	realParams := make([]reflect.Value, 0, len(params))
	for _, param := range params {
		realParams = append(realParams, reflect.ValueOf(param))
	}

	return fv.Call(realParams)
}

func Invoke0(params ...interface{}) []reflect.Value {
	if len(params) == 0 {
		return nil
	}

	if reflect.TypeOf(params[0]).Kind() != reflect.Func {
		return nil
	}

	fv := reflect.ValueOf(params[0])
	realParams := make([]reflect.Value, 0, len(params)-1)
	for i := 1; i < len(params); i++ {
		realParams = append(realParams, reflect.ValueOf(params[i]))
	}

	return fv.Call(realParams)
}
