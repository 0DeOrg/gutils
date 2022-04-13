package mathutils

/**
 * @Author: lee
 * @Description:
 * @File: common
 * @Date: 2022-04-13 4:45 下午
 */
import (
	"github.com/shopspring/decimal"
	"math"
)

// KeepSignificantDigits
/* @Description:保留有效数字
 * @param d decimal.Decimal
 * @param digit int 保留多少位
 * @return decimal.Decimal
 */
func KeepSignificantDigits(d decimal.Decimal, digit int) decimal.Decimal {
	if digit <= 0 {
		return d
	}

	coe := d.Coefficient()
	coeStr := d.Coefficient().String()
	if len(coeStr) > digit {
		pow := decimal.NewFromFloat(math.Pow10(len(coeStr) - digit))
		coe = decimal.NewFromBigInt(d.Coefficient(), 0).DivRound(pow, 0).Mul(pow).BigInt()
	}

	ret := decimal.NewFromBigInt(coe, d.Exponent())

	return ret
}
