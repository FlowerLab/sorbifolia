package httputils

import (
	"github.com/valyala/fasthttp"
)

func (h *HTTP) AddHeader(k, v string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.Add(k, v)
		return nil
	})
}

func (h *HTTP) SetHeader(k, v string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.Set(k, v)
		return nil
	})
}

func (h *HTTP) DelHeader(k string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.Del(k)
		return nil
	})
}

func (h *HTTP) SetReferer(referer string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetReferer(referer)
		return nil
	})
}

func (h *HTTP) SetUserAgent(userAgent string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetUserAgent(userAgent)
		return nil
	})
}

func (h *HTTP) SetProtocol(p string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetProtocol(p)
		return nil
	})
}

func (h *HTTP) SetByteRange(startPos, endPos int) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetByteRange(startPos, endPos)
		return nil
	})
}

func (h *HTTP) SetContentLength(contentLength int) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetContentLength(contentLength)
		return nil
	})
}

func (h *HTTP) SetContentEncoding(contentEncoding string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetContentEncoding(contentEncoding)
		return nil
	})
}

func (h *HTTP) SetMultipartFormBoundary(boundary string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetMultipartFormBoundary(boundary)
		return nil
	})
}
