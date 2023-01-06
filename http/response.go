package http

import (
	"bytes"
	"errors"
	"io"
	"strconv"

	"go.x2ox.com/sorbifolia/http/httpbody"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/kv"
	"go.x2ox.com/sorbifolia/http/render"
	"go.x2ox.com/sorbifolia/http/status"
	"go.x2ox.com/sorbifolia/http/version"
)

type Response struct {
	StatusCode status.Status
	Header     httpheader.ResponseHeader
	Body       io.Reader
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
	r.Header.ContentType = rend.ContentType()
	r.Header.ContentLength = strconv.AppendInt(r.Header.ContentLength, rend.Length(), 10) // need to try to optimize
}

func (r *Response) Encode(ver version.Version) (io.ReadCloser, error) {
	if r.Body != nil && r.Header.ContentLength.Length() == 0 {
		if bytes.Equal(r.Header.Get([]byte("Transfer-Encoding")).V, char.Chunked) {
			panic("TODO: support chunked encoding")
		}
		return nil, errors.New("ContentLength must set")
	}

	var (
		body = r.Body
		buf  = httpbody.AcquireMemory()
	)

	buf.Write(ver.Bytes())
	buf.Write(char.Spaces)
	buf.Write(r.StatusCode.Bytes())
	buf.Write(char.CRLF)

	r.Header.Add(kv.KV{
		K: char.ContentLength,
		V: r.Header.ContentLength,
	})
	r.Header.Each(func(kv kv.KV) bool {
		switch {
		case bytes.EqualFold(kv.K, char.ContentLength):
			if len(r.Header.ContentLength) != 0 {
				kv.V = r.Header.ContentLength
			}
		case bytes.EqualFold(kv.K, char.ContentType):
			if len(r.Header.ContentType) != 0 {
				kv.V = r.Header.ContentType
			}
		case bytes.EqualFold(kv.K, char.SetCookie):
			if len(r.Header.SetCookies) != 0 {
				// kv.V = (*[]byte)(&r.Header.SetCookie)
			}
		case bytes.EqualFold(kv.K, char.Connection):
		}

		buf.Write(kv.K)
		buf.Write(char.Colons)
		buf.Write(kv.V)
		buf.Write(char.CRLF)

		return true
	})

	buf.Write(char.CRLF)
	rio := &responseIO{}
	rio.r = make([]io.Reader, 1, 2)
	buf.Close()
	rio.r[0] = buf
	if body != nil {
		rio.r = append(rio.r, body)
		if c, ok := body.(io.Closer); ok {
			rio.c = make([]io.Closer, 1, 1)
			rio.c[0] = c
		}
	}

	return rio, nil
}

type responseIO struct {
	c []io.Closer
	r []io.Reader
}

func (r *responseIO) Read(p []byte) (int, error) {
	if len(r.r) == 0 {
		return 0, io.EOF
	}

	idx := 0
	for {
		if len(r.r) == 0 || idx == len(p) {
			return idx, nil
		}

		n, err := r.r[0].Read(p[idx:])
		idx += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.r = r.r[1:]
			} else {
				return idx, err
			}
		}
	}
}

func (r *responseIO) Close() (err error) {
	for _, v := range r.c {
		if bErr := v.Close(); err == nil {
			err = bErr
		}
		if p, ok := v.(httpbody.Pool); ok {
			httpbody.Release(p)
		}
	}
	return
}
