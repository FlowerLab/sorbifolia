package httputils

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

type HTTP struct {
	fn     []Handler
	client *fasthttp.Client
	buf    *bytes.Buffer
}

func (h *HTTP) Add(fn ...Handler) *HTTP { h.fn = append(h.fn, fn...); return h }

type Handler func(
	client *fasthttp.Client,
	req *fasthttp.Request,
	resp *fasthttp.Response,
) error
