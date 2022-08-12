package httputils

import (
	"github.com/valyala/fasthttp"
)

func (h *HTTP) SetCookie(k, v string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetCookie(k, v)
		return nil
	})
}

func (h *HTTP) DelCookie(k string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.DelCookie(k)
		return nil
	})
}

func (h *HTTP) DelAllCookies() *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.DelAllCookies()
		return nil
	})
}
