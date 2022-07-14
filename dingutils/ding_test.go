package dingutils

import "testing"

/**
 * @Author: lee
 * @Description:
 * @File: ding_test
 * @Date: 2022-07-14 3:14 下午
 */

func Test_Ding(t *testing.T) {
	InitDingBot("https://oapi.dingtalk.com/robot/send?access_token=4eef40fa5bfa757f278aa66c6083d3f5a266e60b41ab9c42684ce58dff189df6")

	param := map[string]interface{}{
		"test": "test",
	}
	PostDingInfo(123, "test", param)
}
