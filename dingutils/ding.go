package dingutils

/**
 * @Author: lee
 * @Description:
 * @File: ding
 * @Date: 2022/1/11 11:52 上午
 */
import (
	"encoding/json"
	"fmt"
	"gitlab.qihangxingchen.com/qt/gutils/logutils"
	"gitlab.qihangxingchen.com/qt/gutils/network"
	"go.uber.org/zap"
	"strings"
	"time"
)

var dingBot *DingTalk

type DingTalk struct {
	Url    string
	client *network.RestAgent
	mapMsg map[int64]string
}

// 初始化告警结构

func InitDingBot(url string) error {
	if url == "" {
		logutils.Warn("InitDingBot url is empty")
		return fmt.Errorf("InitDingBot url is empty")
	}
	client, err := network.NewRestClient(url, 0, false)
	if nil != err {
		logutils.Fatal("InitDingBot url is unreachable", zap.String("url", url))
	}
	dingBot = &DingTalk{
		Url:    url,
		client: client,
	}
	return nil
}

func AssignMsgMap(mapMsg map[int64]string) {
	if nil != dingBot {
		dingBot.mapMsg = mapMsg
	}
}

func PostDingInfo(code int64, funcName string, params map[string]interface{}) error {
	if nil == dingBot {
		return fmt.Errorf("PostDingInfo dingBot is nil")
	}
	mapParam := params
	if nil == mapParam {
		mapParam = make(map[string]interface{})
	}
	mapParam["code"] = code
	mapParam["func"] = funcName

	return doPostDingMsg(KindInfo, code, mapParam)
}

func PostDingWarn(code int64, funcName string, params map[string]interface{}) error {
	if nil == dingBot {
		return fmt.Errorf("PostDingWarn dingBot is nil")
	}
	mapParam := params
	if nil == mapParam {
		mapParam = make(map[string]interface{})
	}
	mapParam["code"] = code
	mapParam["func"] = funcName

	return doPostDingMsg(KindWarn, code, mapParam)
}

func PostDingError(code int64, funcName string, params map[string]interface{}) error {
	if nil == dingBot {
		return fmt.Errorf("PostDingError dingBot is nil")
	}
	mapParam := params
	if nil == mapParam {
		mapParam = make(map[string]interface{})
	}
	mapParam["code"] = code
	mapParam["func"] = funcName

	return doPostDingMsg(KindError, code, mapParam)
}

func doPostDingMsg(kind string, code int64, params map[string]interface{}) error {
	var title = "未定义告警"
	var ok bool

	str, ok := DING_WARNING_MSG[code]
	if ok {
		title = str
	}

	if nil != dingBot.mapMsg {
		str, ok = dingBot.mapMsg[code]
		if ok {
			title = str
		}
	}

	msg, err := json.Marshal(params)
	if nil != err {
		return err
	}

	sections := strings.Split(string(msg), ",")

	var content strings.Builder
	titleStr := fmt.Sprintf("#### 【%s】", kind)
	content.WriteString(titleStr + title + "\n")

	for _, value := range sections {
		content.WriteString("> " + value + "\n")
	}

	content.WriteString("\n\n ###### ")
	content.WriteString(time.Now().Format("2006-01-02 15:04:05") + "（UTC+8）")

	_, err = postDing(content.String())
	if err != nil {
		return fmt.Errorf("doPostDingMsg|Ding err %s", err.Error())
	}

	return nil
}

// 发送告警信息
func postDing(text string) (string, error) {
	// 构造告警请求
	markDownInfo := &WarnDataReq{
		Title: "告警信息",
		Text:  text,
	}
	atInfo := &WarnAtReq{
		IsAtAll: false,
	}
	warnReq := &WarnReq{
		MsgType:  "markdown",
		MarkDown: *markDownInfo,
		At:       *atInfo,
	}

	reqBody, _ := json.Marshal(warnReq)

	// 发送告警信息
	body, err := dingBot.client.SimplePost("", string(reqBody), nil)

	//body, err := network.HttpPostJson(dingBot.Url, string(reqBody))
	if err != nil {
		return "", err
	}

	return body, nil
}
