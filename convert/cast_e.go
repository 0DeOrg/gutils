package convert

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
)

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

func StructToMapStringE(obj interface{}) (map[string]string, error) {
	err := MustBeStruct(obj)
	if nil != err {
		return nil, err
	}

	reqBody, err := json.Marshal(obj)
	mapParam := make(map[string]interface{})
	err = json.Unmarshal(reqBody, &mapParam)
	if nil != err {
		return nil, err
	}

	ret, err := ToStringMapStringE(mapParam)
	if nil != err {
		return nil, err
	}

	return ret, nil
}

func SliceToStruct(list []interface{}, obj interface{}) error {
	err := MustBeStructPtr(obj)
	if nil != err {
		return err
	}

	objValue := reflect.ValueOf(obj).Elem()
	objType := reflect.TypeOf(obj).Elem()
	for i := 0; i < len(list); i++ {
		if i >= objValue.NumField() {
			break
		}

		src := list[i]
		dst := objValue.Field(i)
		dstType := objType.Field(i)

		if reflect.TypeOf(src).Kind() != dst.Kind() {
			return fmt.Errorf("can't parse '%s' to '%s' for param '%s'", reflect.TypeOf(src).Kind(), dst.Kind(), dstType.Name)
		}

		if dst.IsValid() && dst.CanSet() {
			dst.Set(reflect.ValueOf(src))
		} else {
			return fmt.Errorf("struct param %s can't be set", dstType.Name)
		}
	}

	return nil
}

func MustBeStructPtr(i interface{}) error {
	dstType := reflect.TypeOf(i)
	if reflect.Ptr != dstType.Kind() || reflect.Struct != dstType.Elem().Kind() {
		return fmt.Errorf("%s must be a struct ptr", dstType.Name())
	}

	return nil
}

func MustBeStruct(i interface{}) error {
	dstType := reflect.TypeOf(i)
	if reflect.Ptr == dstType.Kind() {
		if reflect.Struct != dstType.Elem().Kind() {
			return fmt.Errorf("'%s' ptr must be a struct", dstType.Name())
		}
	} else if reflect.Struct != dstType.Kind() {
		return fmt.Errorf("'%s' must be a struct", dstType.Name())
	}

	return nil
}
