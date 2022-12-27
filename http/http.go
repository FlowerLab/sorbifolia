//go:build goexperiment.arenas

package http

import (
	"io"

	"go.x2ox.com/sorbifolia/http/render"
)

type Response struct {
	Header ResponseHeader
	Body   io.Reader
}

func (r *Response) SetBody(body any) {
	if body == nil {
		return
	}

	var rend render.Render
	switch body := body.(type) {
	case string:
		rend = render.Text(body)
	case []byte:
		rend = render.Text(body)
	case render.Render:
		rend = body
	}

	r.Body = rend.Render()
	r.Header.ContentType = rend.ContentType()
	r.Header.ContentLength = rend.Length()
}
