package dingutils

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2022/1/11 11:52 上午
 */

// 告警请求
type WarnReq struct {
	MsgType  string      `json:"msgtype"     binding:"required"` // 消息类型
	MarkDown WarnDataReq `json:"markdown"    binding:"required"` // 告警数据
	At       WarnAtReq   `json:"at"          binding:"required"` // 命中
}

// 告警数据内容
type WarnDataReq struct {
	Title string `json:"title"     binding:"required"` // 主题
	Text  string `json:"text"     binding:"required"`  // 详细内容
}

// 告警数据内容
type WarnAtReq struct {
	IsAtAll bool `json:"isAtAll"     binding:"required"` // 是否全部命中
}
