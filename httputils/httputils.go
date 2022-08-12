package httputils

import (
	"time"

	"github.com/valyala/fasthttp"
)

func (h *HTTP) Request(retry int, isFail func(error, *fasthttp.Response) bool, timeout time.Duration) *HTTP {
	if retry <= 0 {
		retry = 1
	}
	if isFail == nil {
		isFail = func(err error, response *fasthttp.Response) bool {
			return err != nil
		}
	}

	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		var err error
		for i := 0; i < retry; i++ {
			if err = client.DoTimeout(req, resp, timeout); !isFail(err, resp) {
				break
			}
		}
		return err
	})
}

func (h *HTTP) DoRelease() error {
	if h.client == nil {
		h.client = &DefaultClient
	}

	req := getRequestBuffer()
	defer req.Put()
	defer h.Put()

	for _, v := range h.fn {
		if err := v(h.client, req.req, req.resp); err != nil {
			return err
		}
	}
	return nil
}

func newUtil(method Method, param ...string) *HTTP {
	http := getHttpBuffer().SetMethod(method)
	if len(param) > 0 {
		http.SetRequestURI(param[0])
	}
	return http
}
