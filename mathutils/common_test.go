package mathutils

import (
	"github.com/shopspring/decimal"
	"testing"
)

/**
 * @Author: lee
 * @Description:
 * @File: common_test
 * @Date: 2022-04-13 5:25 下午
 */

func Test_Significant(t *testing.T) {

	val := 0.0000123456879

	println(KeepSignificantDigits(decimal.NewFromFloat(val), 6).String())

	val = 0.1234
	println(KeepSignificantDigits(decimal.NewFromFloat(val), 6).String())

	val = 0.1534
	println(KeepSignificantDigits(decimal.NewFromFloat(val), 1).String())

	val = 56789.1234
	println(KeepSignificantDigits(decimal.NewFromFloat(val), 6).String())

	val = 56700.0200
	println(KeepSignificantDigits(decimal.NewFromFloat(val), 6).String())

	val = 56781.0001
	println(KeepSignificantDigits(decimal.NewFromFloat(val), 3).String())
}
