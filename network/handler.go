package network

import "context"

/**
 * @Author: lee
 * @Description:
 * @File: handler
 * @Date: 2022-10-18 3:31 下午
 */

type MsgHandlerCallback func(client *WebsocketAgent, msg string)
type Handler struct {
	client   *WebsocketAgent
	msgIn    chan string
	callback MsgHandlerCallback
	quit     context.Context
}

func NewHandler(client *WebsocketAgent, callback MsgHandlerCallback) *Handler {
	ret := &Handler{
		client:   client,
		msgIn:    make(chan string, 1000),
		callback: callback,
	}

	return ret
}

func (h *Handler) Deliver(msg string) {
	h.msgIn <- msg
}

func (h *Handler) Run() {
	var msg string
	for {
		select {
		case msg = <-h.msgIn:
			h.callback(h.client, msg)
		}
	}
}
