package network
/**
 * @Author: lee
 * @Description:
 * @File: websocket_client
 * @Date: 2021/9/9 11:24 上午
 */

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gutils/logutils"
	"net/url"
	"strconv"
	"time"
)


const (
	MessagePrefix = "WsPrefix:"
)

type  WebsocketAgent struct {
	NetAgentBase
	client *websocket.Conn
	reqChan chan string
	OnPing func(string) error
	OnMessage func(*WebsocketAgent, string) //收到消息回调
	OnSend func(*WebsocketAgent, int, string)    //发送消息回调
	OnClose func(*WebsocketAgent)
	OnConnected func()                          //连接被断开回调
	errConn error
	sendElapse int				//发送消息时间间隔 单位ms 用于限频
}

func NewWebsocketAgent(host string, port uint, path string, isSecure bool, elapse int) *WebsocketAgent {

	hostUrl := ""
	if isSecure {
		hostUrl = "wss://" + host
	} else {
		hostUrl = "ws://" + host
	}

	if 0 != port {
		hostUrl += ":" + strconv.FormatUint(uint64(port), 10)
	}

	hostUrl += path

	rawUrl, err := url.Parse(hostUrl)
	if nil != err {
		panic(err.Error())
	}

	ret := &WebsocketAgent{
		NetAgentBase: NetAgentBase{
			URL: rawUrl,
			isAlive: false,
			timeout: 20,
			isClosed: false,
		},
		reqChan: make(chan string, 128),
		sendElapse: elapse ,
	}

	return ret
}

func (ws *WebsocketAgent)SetPingHandler(handler func(string) error) {
	ws.client.SetPingHandler(handler)
}

func (ws *WebsocketAgent)SetCloseHandler(handler func(code int, text string)error) {
	ws.client.SetCloseHandler(handler)
}

func (ws *WebsocketAgent) Connect() {
	go func() {
		for {
			if ws.isClosed {
				break
			}

			if !ws.isAlive && !ws.isClosed {
				ws.dial()
			}

			time.Sleep(time.Duration(ws.timeout) * time.Millisecond)
		}
	}()
	ws.doSendThread()
	ws.doReceiveThread()
}

func (ws *WebsocketAgent) Close() error {
	ws.isClosed = true
	if nil != ws.client {
		return ws.client.Close()
	}
	return nil
}

func (ws *WebsocketAgent) Send(msg string) {
	messageType := fmt.Sprintf("%02d", websocket.TextMessage)

	ws.reqChan <- MessagePrefix + messageType + msg
}

func (ws *WebsocketAgent) SendPongMsg(data []byte) {
	messageType := fmt.Sprintf("%02d", websocket.PongMessage)
	ws.reqChan <- MessagePrefix + messageType + string(data)
}

func (ws *WebsocketAgent) WaitForConnected() error {
	var ret error
	tick := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <- tick: {
			if ws.isAlive {
				return nil
			}
		}
		case <-time.After(30 * time.Second): {
			ret = fmt.Errorf("wait for websocket connect time out, url: " + ws.URL.String())
			return ret
		}
		}
	}
}

func (ws *WebsocketAgent) dial() error {
	var err error
	urlStr := ws.URL.String()
	logutils.Warn("dial websocket", zap.String("url", ws.URL.String()))
	ws.client, _, err = websocket.DefaultDialer.Dial(urlStr, nil)

	if nil != err {
		ws.errConn = err
		return err
	}

	if nil != ws.OnPing {
		ws.client.SetPingHandler(ws.OnPing)
	}

	ws.client.SetCloseHandler(func(code int, text string) error {
		ws.isAlive = false
		if nil !=  ws.OnClose {
			ws.OnClose(ws)
		}
		message := websocket.FormatCloseMessage(code, "")
		ws.client.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	ws.isAlive = true
	ws.OnConnected()
	ws.errConn = nil

	return nil
}

func (ws *WebsocketAgent) doSendThread() {
	//logutils.Warn("doSendThread", zap.String("url", ws.URL.String()))
	go func(){
		for {
			if ws.isClosed {
				break
			}

			msg := <-ws.reqChan

			messageType, sendMsg := ParseMessage(msg)
			var err error
			switch messageType {
			case websocket.TextMessage, websocket.BinaryMessage:
				err = ws.client.WriteMessage(messageType, []byte(sendMsg))
				break
			case websocket.PongMessage, websocket.PingMessage, websocket.CloseMessage:
				err = ws.client.WriteControl(messageType, []byte(sendMsg), time.Now().Add(time.Second))
				break
			}

			if nil != err {
				ws.isAlive = false
				logutils.Warn("doSendThread fatal", zap.String("url", ws.URL.String()), zap.Error(err))
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if nil != ws.OnSend {
				ws.OnSend(ws, messageType, sendMsg)
			}
			if ws.sendElapse > 0 {
				elapse := time.Duration(ws.sendElapse) * time.Millisecond
				time.Sleep( elapse )
			}
		}
	}()
}

func (ws *WebsocketAgent) doReceiveThread() {
	//logutils.Warn("doReceiveThread", zap.String("url", ws.URL.String()))
	go func() {
		for {
			if ws.isClosed{
				break
			}
			if !ws.isAlive {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			_, msg, err := ws.client.ReadMessage()
			if nil != err {
				ws.isAlive = false
				logutils.Warn("doReceiveThread fatal", zap.String("url", ws.URL.String()), zap.Error(err))
				continue
			}

			ws.OnMessage(ws, string(msg))
		}
	}()
}

func ParseMessage(msg string) (int, string) {
	prefix := msg[0: len(MessagePrefix)]
	if prefix != MessagePrefix {
		return 0, ""
	}

	remain := msg[len(MessagePrefix):]
	msgType := remain[0: 2]
	messageType, _ := strconv.ParseInt(msgType, 10, 64)
	remain = remain[2:]

	return int(messageType), remain
}
