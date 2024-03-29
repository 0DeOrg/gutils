package network

/**
 * @Author: lee
 * @Description:
 * @File: rest
 * @Date: 2021/9/7 3:48 下午
 */

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RestAgent struct {
	NetAgentBase
	Client *resty.Client
}

var _ HttpInterface = (*RestAgent)(nil)

func NewRestClient(host string, port uint, isHttps bool) (*RestAgent, error) {
	hostUrl := ""

	trimHost := strings.TrimLeft(host, " ")

	if strings.HasPrefix(trimHost, "http") && strings.Contains(trimHost, "://") {
		hostUrl = trimHost
	} else {
		if isHttps {
			hostUrl += "https://" + trimHost
		} else {
			hostUrl += "http://" + trimHost
		}
	}

	if 0 != port {
		hostUrl += ":" + strconv.FormatUint((uint64)(port), 10)
	}

	url, err := url.Parse(hostUrl)
	if nil != err {
		return nil, err
	}

	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	hc := http.Client{
		Jar:     cookieJar,
		Timeout: 20 * time.Second,
	}
	client := resty.NewWithClient(&hc)
	ret := RestAgent{
		NetAgentBase: NetAgentBase{
			URL: url,
		},
		Client: client,
	}

	return &ret, nil
}

func (h *RestAgent) SimpleGet(path string, params map[string]string) (string, error) {
	url := h.URL.String() + path
	if nil != params {

	}
	res, err := h.Client.R().SetQueryParams(params).Get(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != 200 {
		return string(res.Body()), fmt.Errorf("response err: %s", res.String())
	}

	return res.String(), nil
}

func (h *RestAgent) SimplePost(path string, reqBody string, params map[string]string) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()

	res, err := r.SetQueryParams(params).SetBody(reqBody).SetHeader("Content-Type", "application/json").Post(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}

	return res.String(), nil
}

func (h *RestAgent) Get(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()
	res, err := r.SetQueryParams(params).SetHeaders(headers).SetCookies(cookies).Get(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}

	return res.String(), nil
}

func (h *RestAgent) Post(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	r := h.Client.R()
	res, err := r.SetQueryParams(params).SetBody(reqBody).SetHeaders(headers).SetCookies(cookies).Post(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil

}

func (h *RestAgent) PostForm(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	res, err := h.Client.R().SetFormData(params).SetBody(reqBody).SetHeaders(headers).SetCookies(cookies).Post(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil
}

func (h *RestAgent) Put(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	res, err := h.Client.R().SetQueryParams(params).SetBody(reqBody).SetHeaders(headers).SetCookies(cookies).Put(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil

}

func (h *RestAgent) Delete(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	res, err := h.Client.R().SetQueryParams(params).SetBody(reqBody).SetHeaders(headers).SetCookies(cookies).Delete(url)
	if nil != err {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response err: %s", res.String())
	}
	return res.String(), nil

}
