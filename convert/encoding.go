package convert

/**
 * @Author: lee
 * @Description:
 * @File: encoding
 * @Date: 2021/9/3 11:20 上午
 */

import "github.com/axgle/mahonia"

func ConvertUnicode2GBK(text string) string {
	enc := mahonia.NewEncoder("gbk")
	ret := enc.ConvertString(text)
	return ret
}