package httputils

import (
	"bufio"
	"errors"
	"io"
	"net/url"

	"github.com/valyala/fasthttp"
)

func (h *HTTP) SetBodyWithEncoder(ge GetEncoder, body any) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		if ge != nil {
			h.buf.Reset()
			err := ge(h.buf).Encode(body)
			if err == nil {
				req.SetBody(h.buf.Bytes())
			}
			return err
		}

		switch body := body.(type) {
		case []byte:
			req.SetBody(body)
		case string:
			req.SetBodyString(body)
		case url.Values:
			req.SetBodyString(body.Encode())
		case nil:
			return nil
		}
		return errors.New("invalid body type")
	})
}

func (h *HTTP) AppendBody(p []byte) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.AppendBody(p)
		return nil
	})
}

func (h *HTTP) AppendBodyString(s string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.AppendBodyString(s)
		return nil
	})
}

func (h *HTTP) ReadBody(r *bufio.Reader, contentLength int, maxBodySize int) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		return req.ReadBody(r, contentLength, maxBodySize)
	})
}

func (h *HTTP) SetBody(body []byte) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.SetBody(body)
		return nil
	})
}

func (h *HTTP) SetBodyStream(bodyStream io.Reader, bodySize int) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.SetBodyStream(bodyStream, bodySize)
		return nil
	})
}

func (h *HTTP) SetBodyString(body string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.SetBodyString(body)
		return nil
	})
}

func (h *HTTP) SetConnectionClose() *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.SetConnectionClose()
		return nil
	})
}

func (h *HTTP) SetRequestURI(requestURI string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.SetRequestURI(requestURI)
		return nil
	})
}

func (h *HTTP) SetHost(host string) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.SetHost(host)
		return nil
	})
}
