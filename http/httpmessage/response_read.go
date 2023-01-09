package httpmessage

import (
	"io"

	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
)

var _ io.ReadCloser = (*Response)(nil)

func (r *Response) Read(p []byte) (n int, err error) {
	if err = r.preprocessRead(); err != nil {
		return
	}

	var wn int

	for n < len(p) {
		switch r.state.Operate() {
		case _Version:
			wn, err = r.writeVersion(p[n:])
		case _Status:
			wn, err = r.writeStatus(p[n:])
		case _Header:
			wn, err = r.writeHeader(p[n:])
		case _Body:
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

func (r *Response) writeVersion(p []byte) (n int, err error) {
	n = copy(p, r.buf.B)
	r.buf.Discard(0, n)

	if r.buf.Len() == 0 {
		r.state.SetOperate(_Status)
		_, _ = r.buf.Write(r.StatusCode.Bytes())
		_, _ = r.buf.Write(char.CRLF)
	}

	return
}

func (r *Response) writeStatus(p []byte) (n int, err error) {
	n = copy(p, r.buf.B)
	r.buf.Discard(0, n)

	if r.buf.Len() == 0 {
		r.state.SetOperate(_Header)
	}

	return
}

func (r *Response) writeHeader(p []byte) (n int, err error) {
	var headerLen = r.Header.KVs.Len()

	for r.p <= headerLen {
		if len(p) == 0 {
			break
		}
		if r.buf.Len() > r.buf.P {
			wn, _ := r.buf.Read(p[n:])
			n += wn
			continue
		}
		r.buf.Reset()

		if r.p == headerLen {
			r.state.SetOperate(_Body)
			break
		}

		r.buf.B = httpheader.AppendHeader(r.buf.B, r.Header.KVs[r.p])
		r.p++
		if r.p == headerLen {
			_, _ = r.buf.Write(char.CRLF)
		}
	}

	return
}

func (r *Response) writeBody(p []byte) (n int, err error) {
	if r.Body == nil {
		return 0, io.EOF
	}
	if n, err = r.Body.Read(p); err == nil {
		return
	}
	r.state.Close()

	if c, ok := r.Body.(io.Closer); ok {
		if cErr := c.Close(); err != nil && err != io.EOF {
			err = cErr
		}
	}
	return
}

func (r *Response) preprocessRead() (err error) {
	if r.state == _Init {
		r.state.SetRead()
		r.state.SetOperate(_Version)

		r.buf = &bufpool.ReadBuffer{}
		_, _ = r.buf.Write(r.Version.Bytes())
		_, _ = r.buf.Write(char.Spaces)

		if r.Body != nil {
			if r.Header.ContentLength().Length() <= 0 {
				return httperr.BodyLengthMismatch
			}
		} else {
			r.Header.SetContentLength(0)
		}
	}

	if !r.state.Readable() {
		return io.EOF
	}

	return nil
}
