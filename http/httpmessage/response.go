package httpmessage

import (
	"io"
	"strconv"

	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/render"
	"go.x2ox.com/sorbifolia/http/status"
	"go.x2ox.com/sorbifolia/http/version"
)

type Response struct {
	Version    version.Version
	StatusCode status.Status
	Header     httpheader.ResponseHeader
	Body       io.Reader

	buf   *bufpool.ReadBuffer
	state state
	p     int
}

func (r *Response) SetBody(body any) {
	if body == nil {
		return
	}

	var rend render.Render
	switch body := body.(type) {
	case string:
		if len(body) == 0 {
			return
		}
		rend = render.Text(body)
	case []byte:
		if len(body) == 0 {
			return
		}
		rend = render.Text(body)
	case render.Render:
		rend = body
	}

	r.Body = rend.Render()
	r.Header.Set([]byte("Content-Type"), rend.ContentType())
	r.Header.Set([]byte("Content-Length"), strconv.AppendInt(nil, rend.Length(), 10))
}

func (r *Response) Close() error {
	switch r.state {
	case 1, 2: // W
		r.state = 0
		r.buf = nil
	case 0, 3: // C
		return nil
	}

	return nil
}

func (r *Response) Reset() {
	r.Header.Reset()
	if r.Body != nil {
		r.Body = nil
	}

	r.buf = nil
	r.state = _Init
	r.p = 0
}
