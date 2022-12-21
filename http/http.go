//go:build goexperiment.arenas

package http

import (
	"bytes"
	"io"

	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type Response struct {
	Header ResponseHeader
	Body   io.Reader
}

func (r *Response) SetBody(body any) {
	switch body := body.(type) {
	case []byte:
		r.Body = bytes.NewReader(body)
	case string:
		r.Body = bytes.NewReader(pyrokinesis.String.ToBytes(body))
	case Render:
		r.Body = body.Render()
		r.Header.ContentType = body.ContentType()
	case io.Reader:
		r.Body = body
	default:
		panic("unknown body")
	}
}
