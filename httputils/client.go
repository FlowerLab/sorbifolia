package httputils

import (
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// DefaultClient is a default client
var DefaultClient = fasthttp.Client{
	MaxConnsPerHost:     50000,
	MaxIdleConnDuration: time.Second * 15, // 空闲 TCP 连接的最大存活时间
	MaxConnDuration:     time.Second * 10, // 空闲 Handler 连接的最大存活时间
	ReadTimeout:         time.Second * 90, // 读取超时
	MaxConnWaitTimeout:  time.Second * 5,  // 等待连接的时间
}

func (h *HTTP) SetClient(client *fasthttp.Client) *HTTP { h.client = client; return h }
func (h *HTTP) SetProxy(addr string) *HTTP {
	if h.client == nil {
		h.client = &fasthttp.Client{}
	}

	_url, _ := url.Parse(addr)
	switch {
	case _url == nil:
	case _url.Scheme == "http", _url.Scheme == "https":
		_url.Scheme = ""
		h.client.Dial = fasthttpproxy.FasthttpHTTPDialer(_url.String())
	case _url.Scheme == "socks5", _url.Scheme == "socks5h":
		h.client.Dial = fasthttpproxy.FasthttpSocksDialer(addr)
	}
	return h
}
