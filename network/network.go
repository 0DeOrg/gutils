package network
/**
 * @Author: lee
 * @Description:
 * @File: network
 * @Date: 2021/9/3 2:24 下午
 */

import (
	"net"
	"net/http"
	"net/url"
)

type HttpInterface interface {
	SimpleGet(path string, params map[string]string) (string, error)
	SimplePost(path string, body string, params map[string]string) (string, error)

	Get(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error)
	Post(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error)
}

type SocketInterface interface {
	Connect()
	Send(string)
	Close()
}

type NetAgentBase struct {
	URL *url.URL
	isAlive 	bool
	timeout 	int
	isClosed	bool
}

var _ HttpInterface = (*NetAgentBase)(nil)

func (b *NetAgentBase) SimpleGet(path string, params map[string]string) (string, error) {
	return "", nil
}

func (b *NetAgentBase) SimplePost(path string, body string, params map[string]string) (string, error) {
	return "", nil
}

func (b *NetAgentBase) Get(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	return "", nil
}

func (b *NetAgentBase) Post(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	return "", nil
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}