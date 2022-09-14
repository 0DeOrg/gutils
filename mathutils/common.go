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
	"math/big"
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

// GetMaxDecimalValue
/* @Description: 获取最高位10进制数并返回位数
 * @param value *big.Int
 * @return *big.Int
 * @return int32
 */
func GetMaxDecimalValue(value *big.Int) (*big.Int, int32) {
	if value.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), 0
	}
	neg := false
	if value.Cmp(big.NewInt(0)) < 0 {
		value = value.Mul(value, big.NewInt(-1))
		neg = true
	}

	str := value.String()

	var digit = int32(len(str) - 1)

	if neg {
		return decimal.NewFromInt32(10).Pow(decimal.NewFromInt32(digit)).Mul(decimal.NewFromInt32(-1)).BigInt(), digit
	}

	return decimal.NewFromInt32(10).Pow(decimal.NewFromInt32(digit)).BigInt(), digit
}
