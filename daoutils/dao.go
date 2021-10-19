package daoutils

/**
 * @Author: lee
 * @Description:
 * @File: dao
 * @Date: 2021/9/15 3:39 下午
 */

import "fmt"


type IDaoClient interface {
	Connect() error
	DSN() string
}

var (
	ErrorNilGormDatabase = fmt.Errorf("gorm database is nil")
)