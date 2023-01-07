package httpmessage

import (
	"io"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

type Request struct {
	cfg        *httpconfig.Config
	state      parseRequestState
	buf        bufpool.Buffer
	rp         int
	bodyLength int

	ver    version.Version
	Method method.Method
	Header httpheader.RequestHeader
	Body   io.ReadWriteCloser
}

func (r *Request) Reset() {
	r.state = 0
	r.buf.Reset()
	r.rp = 0
	r.bodyLength = 0
	r.Header.Reset()
	if r.Body != nil {
		_ = r.Body.Close()
		r.Body = nil
	}
}
