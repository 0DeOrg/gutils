package network

/**
 * @Author: lee
 * @Description:
 * @File: signalR
 * @Date: 2022/1/18 10:18 上午
 */

import (
	"context"
	kitlog "github.com/go-kit/log"
	"github.com/philippseith/signalr"
	"net/url"
	"os"
	"strconv"
)

type agentHub struct {
	signalr.Hub
}

func (h *agentHub) OnConnected(connectionID string) {

}

func (h *agentHub) OnDisconnected(connectionID string) {

}

type SignalRAgent struct {
	NetAgentBase
	client signalr.Client
	ctx    context.Context
}

func NewSignalRAgent(host string, path string, port uint, hubs []string, isSecure bool, receiver signalr.ReceiverInterface) (*SignalRAgent, error) {
	hostUrl := ""
	if isSecure {
		hostUrl = "https://" + host
	} else {
		hostUrl = "http://" + host
	}

	if 0 != port {
		hostUrl += ":" + strconv.FormatUint(uint64(port), 10)
	}

	hostUrl += path

	rawUrl, err := url.Parse(hostUrl)
	if nil != err {
		panic(err.Error())
	}

	client, err := signalr.NewClient(context.Background(),
		signalr.WithReceiver(receiver),
		signalr.WithAutoReconnect(func() (signalr.Connection, error) {
			return signalr.NewHTTPConnection(context.TODO(), hostUrl)
		}),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stdout), false))
	if nil != err {
		return nil, err
	}

	ret := &SignalRAgent{
		NetAgentBase: NetAgentBase{
			URL: rawUrl,
		},
		client: client,
	}

	return ret, nil
}

func (sig *SignalRAgent) Connect() {
	sig.client.Start()
}

func (sig *SignalRAgent) WaitForConnected() error {
	return <-sig.client.WaitForState(context.Background(), signalr.ClientConnected)
}

func (sig *SignalRAgent) Invoke(method string, arguments ...interface{}) signalr.InvokeResult {
	return <-sig.client.Invoke(method, arguments...)
}

func (sig *SignalRAgent) Send(method string, arguments ...interface{}) error {
	return <-sig.client.Send(method, arguments...)
}
