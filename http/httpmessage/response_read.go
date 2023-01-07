package httpmessage

import (
	"fmt"
	"io"

	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
)

var _ io.ReadCloser = (*Response)(nil)

func (r *Response) Read(p []byte) (n int, err error) {
	switch r.state {
	case 0: // RW
		r.buf = &bufpool.ReadBuffer{}
		r.state = 1
	case 1: // R
	case 2, 3: // C
		return 0, io.EOF
	}

	var wn int

	for n < len(p) {
		switch r.status {
		case 0: // 0 init，1 status，2 header，3 body, 4 END
			_, _ = r.buf.Write(r.StatusCode.Bytes())
			_, _ = r.buf.Write(char.CRLF)
			r.status = 1
			continue
		case 1:
			wn, err = r.writeStatus(p[n:])
		case 2:
			wn, err = r.writeHeader(p[n:])
		case 3:
			wn, err = r.writeBody(p[n:])
		default:
			panic("?")
		}

		n += wn

		if err != nil {
			break
		}
	}

	return
}

func (r *Response) writeStatus(p []byte) (n int, err error) {
	n = copy(p, r.buf.B[r.p:])
	r.buf.Discard(0, n)

	if r.buf.Len() == 0 {
		r.status++
	}

	return
}

func (r *Response) writeHeader(p []byte) (n int, err error) {
	var headerLen = len(r.Header.KVs)

	for r.p <= headerLen {
		if len(p) == 0 {
			break
		}
		if r.buf.Len() > r.buf.P {
			wn, _ := r.buf.Read(p)
			n += wn
			p = p[wn:]
			continue
		}
		r.buf.Reset()

		if r.p >= headerLen {
			r.status++
			break
		}

		r.buf.B = r.Header.KVs[r.p].AppendHeader(r.buf.B)
		r.p++
		if r.p == headerLen {
			_, _ = r.buf.Write(char.CRLF)
		}
	}

	return
}

func (r *Response) writeBody(p []byte) (n int, err error) {
	if n, err = r.Body.Read(p); err == nil {
		fmt.Println(err)
		return
	}
	r.status++

	if c, ok := r.Body.(io.Closer); ok {
		if cErr := c.Close(); err != nil {
			err = cErr
		}
	}
	return
}
