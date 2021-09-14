package network

import (
	"github.com/btcsuite/websocket"
	"go.uber.org/zap"
	"gutils/logger"
	"net/url"
	"strconv"
	"time"
)

/**
 * @Author: lee
 * @Description:
 * @File: websocket_client
 * @Date: 2021/9/9 11:24 上午
 */

type WebsocketClient struct {
	BaseClient
	client *websocket.Conn
	reqChan chan string
	OnMessage func(string)
}

func NewWebsocketClient(host string, port uint, path string, isSecure bool) (*WebsocketClient, error) {
	hostUrl := ""
	if isSecure {
		hostUrl += "wss://" + host
	} else {
		hostUrl += "ws://" + host
	}

	if 0 != port {
		hostUrl += ":" + strconv.FormatUint(uint64(port), 10)
	}

	hostUrl += path

	rawUrl, err := url.Parse(hostUrl)
	if nil != err {
		return nil, err
	}

	ret := &WebsocketClient{
		BaseClient: BaseClient{
			URL: rawUrl,
			isAlive: false,
			timeout: 20,
			isClosed: false,
		},
		reqChan: make(chan string, 128),
	}

	return ret, nil
}

func (ws *WebsocketClient) Connect() {
	go func() {
		for {
			if ws.isClosed {
				break
			}

			if !ws.isAlive && !ws.isClosed {
				ws.dial()
				ws.doSendThread()
				ws.doReceiveThread()
			}

			time.Sleep(time.Duration(ws.timeout) * time.Second)
		}
	}()
}

func (ws *WebsocketClient) Close() {
	ws.isClosed = true
	ws.client.Close()
}

func (ws *WebsocketClient) Send(msg string) {
	ws.reqChan <- msg
	i := 0
	i++
}

func (ws *WebsocketClient) dial() error {
	var err error
	ws.client, _, err = websocket.DefaultDialer.Dial(ws.URL.String(), nil)
	if nil != err {
		return err
	}

	ws.isAlive = true

	return nil
}

func (ws *WebsocketClient) doSendThread() {
	go func(){
		for {
			if ws.isClosed {
				break
			}
			msg := <-ws.reqChan
			logger.Info("websocket send msg", zap.String("url", ws.URL.String()), zap.String("msg", msg))
			err := ws.client.WriteMessage(websocket.TextMessage, []byte(msg))
			if nil != err {
				ws.isAlive = false
				break
			}
		}
	}()
}

func (ws *WebsocketClient) doReceiveThread() {
	go func() {
		for {
			if ws.isClosed {
				break
			}
			_, msg, err := ws.client.ReadMessage()
			if nil != err {
				ws.isAlive = false
				break
			}

			ws.OnMessage(string(msg))
		}
	}()
}

