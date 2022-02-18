package dingutils

/**
 * @Author: lee
 * @Description:
 * @File: constant
 * @Date: 2022/1/11 11:54 上午
 */

const (
	KindError = "错误"
	KindWarn  = "警告"
	KindInfo  = "通知"
)

const (
	WarnIndexPriceMasterSlaveChange = iota + 1000
	WarnIndexPriceSourceOffsetReachAlert
	WarnIndexPriceIsOld
)

var DING_WARNING_MSG = map[int64]string{
	WarnIndexPriceMasterSlaveChange:      "指数价格主备切换",
	WarnIndexPriceSourceOffsetReachAlert: "数据源当前价格时间偏离当前时间到达告警值",
	WarnIndexPriceIsOld:                  "当前指数价格5s没有生成新的",
}
