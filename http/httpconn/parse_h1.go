package httpconn

import (
	"bytes"
	"crypto/tls"
	"io"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/url"
)

type Request struct {
	Method method.Method
	URL    url.URL
	Header *httpheader.Header
	Body   io.ReadWriteCloser

	RequestURI []byte
	TLS        *tls.ConnectionState
}

var (
	handle10 = make(chan HTTPConn, 1)
	handle11 = make(chan HTTPConn, 1)
	handle2  = make(chan HTTPConn, 1)
	handle3  = make(chan HTTPConn, 1)
)

func run2(conn HTTPConn) { _ = conn.Close() }
func run1(conn HTTPConn, buf *bufpool.Buffer, req *Request) {
	req.Header = httpheader.Acquire()

	hw := &HeaderWriterWith1{
		cfg:    nil,
		Header: req.Header,
		buf:    buf,
	}
	n, err := hw.Write(buf.B)
	if err != nil {
		if err != io.EOF {
			return
		}
	} else if _, err = util.Copy(hw, conn); err != nil && err != io.EOF {
		return
	}
	if n < buf.Len() {
		// write body
	}

}

type HeaderWriterWith1 struct {
	cfg    *httpconfig.Config
	Header *httpheader.Header
	buf    *bufpool.Buffer
}

func (r *HeaderWriterWith1) Write(p []byte) (n int, err error) {
	n = len(p)
	b := p

	for len(b) > 0 {
		// Check "\r\n" is straddles the buffer.
		if length := r.buf.Len(); length != 0 && r.buf.B[length-1] == '\r' && p[0] == '\n' {
			r.Header.AddHeader(r.buf.B[:length-1])
			r.buf.Reset()
			b = b[1:]
			continue
		}

		i := bytes.Index(p, char.CRLF)
		switch i {
		case 0:
			if r.buf.Len() == 0 {
				err = io.EOF
				_, _ = r.buf.Write(p[2:])
				return
			}
		case -1:
			if size := r.cfg.GetMaxRequestHeaderSize(); size > 0 && r.buf.Len()+len(b) > size {
				return 0, httperr.RequestHeaderFieldsTooLarge
			}
			_, _ = r.buf.Write(b)
			return
		default:
			if size := r.cfg.GetMaxRequestHeaderSize(); size > 0 && r.buf.Len()+i > size {
				return 0, httperr.RequestHeaderFieldsTooLarge
			}
			_, _ = r.buf.Write(b[:i])
			b = b[i+2:]
		}

		r.Header.AddHeader(r.buf.B)
		r.buf.Reset()
	}

	return
}
