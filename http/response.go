//go:build goexperiment.arenas

package http

import (
	"arena"
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/render"
	"go.x2ox.com/sorbifolia/http/version"
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
	r.Header.ContentLength = strconv.AppendInt(r.Header.ContentLength, rend.Length(), 10) // need to try to optimize
}

func (r *Response) Encode(ver version.Version, a *arena.Arena) (io.ReadCloser, error) {
	if r.Body != nil && r.Header.ContentLength.Length() == 0 {
		if r.Header.Get([]byte("Transfer-Encoding")).Equal(char.Chunked) {
			// TODO: support chunked encoding
		}
		return nil, errors.New("ContentLength must set")
	}

	var (
		body             = r.Body
		nb   net.Buffers = arena.MakeSlice[[]byte](a, 4, 4+len(r.Header.KVs)*4+1)
	)

	{
		nb[0] = ver.Bytes()
		nb[1] = char.Spaces
		nb[2] = r.Header.StatusCode.Bytes()
		nb[3] = char.CRLF
	}
	bufAppend := func(b []byte) {
		if length := len(nb) + 1; length < cap(nb) {
			arr := arena.MakeSlice[[]byte](a, len(nb), length*5/4)
			copy(arr, nb)
			nb = arr
		}
		nb = append(nb, b)
	}

	r.Header.add(KV{
		K: char.ContentLength,
		V: (*[]byte)(&r.Header.ContentLength),
	})
	r.Header.Each(func(kv KV) bool {
		switch {
		case bytes.EqualFold(kv.K, char.ContentLength):
			if len(r.Header.ContentLength) != 0 {
				kv.V = (*[]byte)(&r.Header.ContentLength)
			}
		case bytes.EqualFold(kv.K, char.ContentType):
			if len(r.Header.ContentType) != 0 {
				kv.V = (*[]byte)(&r.Header.ContentType)
			}
		case bytes.EqualFold(kv.K, char.SetCookie):
			if len(r.Header.SetCookies) != 0 {
				// kv.V = (*[]byte)(&r.Header.SetCookie)
			}
		case bytes.EqualFold(kv.K, char.Connection):
		}

		bufAppend(kv.K)
		bufAppend(char.Colons)
		bufAppend(kv.Val())
		bufAppend(char.CRLF)

		return true
	})

	bufAppend(char.CRLF)
	rio := arena.New[responseIO](a)
	rio.r = arena.MakeSlice[io.Reader](a, 1, 2)
	rio.r[0] = &nb
	if body != nil {
		rio.r = append(rio.r, body)
		if c, ok := body.(io.Closer); ok {
			rio.c = arena.MakeSlice[io.Closer](a, 1, 1)
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
	}
	return
}
