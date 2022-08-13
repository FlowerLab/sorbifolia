package httputils

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func (h *HTTP) test() (*fasthttp.Request, *fasthttp.Response, error) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	for _, v := range h.fn {
		if err := v(h.client, req, resp); err != nil {
			return req, resp, err
		}
	}

	return req, resp, nil
}

func TestHTTP_Add(t *testing.T) {
	h := Post().Add(nil)
	if len(h.fn) != 2 {
		t.Error("Add err")
	}
}
