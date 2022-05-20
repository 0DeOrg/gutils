package fileutils

import (
	"os"
)

/**
 * @Author: lee
 * @Description:
 * @File: file
 * @Date: 2021/9/14 2:26 下午
 */

func PathExist(filePath string) bool {
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CreateDirectoryIfNotExist(pathName string, mode os.FileMode) error {
	exist := PathExist(pathName)
	if !exist {
		err := os.MkdirAll(pathName, mode)
		if err != nil {
			return err
		}
	}

	return nil
}
