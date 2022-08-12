package httputils

import (
	"bytes"
	"sync"

	"github.com/valyala/fasthttp"
)

type requestBuffer struct {
	req  *fasthttp.Request
	resp *fasthttp.Response
}

func (r *requestBuffer) Put() {
	if r == nil || r.req == nil || r.resp == nil {
		return
	}
	r.req.Reset()
	r.resp.Reset()
	requestPool.Put(r)
}

var (
	requestPool = &sync.Pool{New: func() interface{} {
		return &requestBuffer{
			req:  &fasthttp.Request{},
			resp: &fasthttp.Response{},
		}
	}}

	httpPool = &sync.Pool{New: func() interface{} {
		return &HTTP{
			buf: &bytes.Buffer{},
		}
	}}
)

func getRequestBuffer() *requestBuffer { return requestPool.Get().(*requestBuffer) }
func getHttpBuffer() *HTTP             { return httpPool.Get().(*HTTP) }

func (h *HTTP) Put() {
	if h == nil || h.buf == nil {
		return
	}
	h.buf.Reset()
	h.fn = h.fn[:0]
	httpPool.Put(h)
}
