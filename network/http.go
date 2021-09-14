package network
/**
 * @Author: lee
 * @Description:
 * @File: http
 * @Date: 2021/9/2 6:07 下午
 */
import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)


type HttpHost struct {
	Host string
	Port int
	IsHttps bool
	CAPath string
}

type HttpClient struct {
	BaseClient
	Client *http.Client
}

//转发消息存储在context中的字段
const (
	ForwardCustomAck = "CustomAck"
	ForwardCustomReq = "CustomReq"
)

func NewHttpClient(host string, port uint, isHttps bool) (*HttpClient, error) {
	hostUrl := ""
	if isHttps {
		hostUrl += "https://" + host
	} else {
		hostUrl += "http://" + host
	}

	if 0 != port {
		hostUrl += ":" + strconv.FormatUint((uint64)(port), 10)
	}

	url, err := url.Parse(hostUrl)
	if nil != err {
		return nil, err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	ret := &HttpClient{
		BaseClient: BaseClient{
			URL: url,
		},
		Client: client,
	}

	return ret, nil
}

var _ HttpInterface = (*HttpClient)(nil)

func (h *HttpClient) SimpleGet(path string, params map[string]string) (string, error) {
	url := h.URL.String() + path
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	if nil != params {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (h *HttpClient) SimplePost(path string, reqBody string, params map[string]string) (string, error){
	url := h.URL.String() + path
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	if nil != params {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(resBody), nil
}

func (h *HttpClient) Get(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	if nil != params {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	for _, v := range cookies {
		req.AddCookie(v)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func (h *HttpClient) Post(path string, reqBody string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (string, error) {
	url := h.URL.String() + path
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if nil != params {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	for _, v := range cookies {
		req.AddCookie(v)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(resBody), nil
}

type transport struct {
	http.RoundTripper
	cxt *gin.Context
}

func (t *transport) RoundTrip (req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if nil != err {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	t.cxt.Set(ForwardCustomAck, string(resBody))

	//err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body := ioutil.NopCloser(bytes.NewReader(resBody))
	resp.Body = body
	resp.ContentLength = int64(len(resBody))
	resp.Header.Set("Content-Length", strconv.Itoa(len(resBody)))


	return resp, err
}

func HttpForward (w http.ResponseWriter, req *http.Request, targetHost *HttpHost, c *gin.Context) error {
	host := ""
	if targetHost.IsHttps {
		host = "https://" + targetHost.Host
	} else {
		host += "http://" + targetHost.Host
	}


	if 0 != targetHost.Port {
		host += ":" + strconv.Itoa(targetHost.Port)
	}

	remote, err := url.Parse(host)
	if nil != err {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if nil != err {
		return err
	}
	body := ioutil.NopCloser(bytes.NewReader(reqBody))
	c.Request.Body = body
	c.Request.ContentLength = int64(len(reqBody))
	c.Request.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

	c.Set(ForwardCustomReq, string(reqBody))

	proxy := httputil.NewSingleHostReverseProxy(remote)
	req.Header.Add("appcode", "app5")
	proxy.Transport = &transport{http.DefaultTransport, c}
	proxy.ServeHTTP(w, req)

	return nil
}
