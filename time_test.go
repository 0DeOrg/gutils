package gutils

import (
	"testing"
	"time"
)

/**
 * @Author: lee
 * @Description:
 * @File: time_test
 * @Date: 2022-04-12 2:53 下午
 */

func Test_Time(t *testing.T) {
	now := time.Now().In(time.UTC).Format(time.RFC3339)
	println(now)
	tm, _ := time.ParseInLocation(time.RFC3339, now, time.UTC)
	println(tm.Format(time.RFC3339))
}
