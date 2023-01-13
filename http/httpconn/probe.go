package httpconn

import (
	"bytes"
	"io"
	"net"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

type Probe struct {
	buf *bufpool.Buffer
	cfg *httpconfig.Config
	fl  FL
}

func (p *Probe) isHTTP2(r io.Reader) (ok bool, err error) {
	if _, err = p.buf.ReadLimit(r, ConnectionPrefaceLength); err == nil {
		ok = bytes.Equal(p.buf.B, ConnectionPreface)
	}
	return
}

type FL struct {
	Method method.Method
	URI    []byte
}

func (p *Probe) ExpandBuffer(size int) {
	if size == 0 {
		size = 64
	}
	if p.buf == nil {
		p.buf = bufpool.AcquireSegment(size)
		return
	}
	if cap(p.buf.B)-p.buf.Len() > size {
		return
	}

	buf := bufpool.AcquireSegment(cap(p.buf.B) * 2)
	_, _ = buf.Write(p.buf.B)
	bufpool.ReleaseSegment(p.buf)
	p.buf = buf
}

func (p *Probe) Probe(conn net.Conn) (HTTPConn, error) {
	p.ExpandBuffer(0)

	if ok, err := p.isHTTP2(conn); err != nil {
		return nil, err
	} else if ok {
		return &Conn{Conn: conn, ver: _h2}, nil
	}

	if err := p.ParseMethod(conn); err != nil {
		return nil, err
	}
	if err := p.ParseURI(conn); err != nil {
		return nil, err
	}
	ver, err := p.ParseVersion(conn)
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: conn, ver: ver}, nil
}

func (p *Probe) ParseURI(r io.Reader) error {
	for {
		i := p.buf.IndexByte(char.Space)
		switch i {
		case 0:
			return httperr.ParseHTTPVersionErr
		case -1:
			if size := p.cfg.GetMaxRequestURISize() - p.buf.Len() - 64; size <= 0 {
				return httperr.RequestURITooLong
			}
			p.ExpandBuffer(64)
			if _, err := p.buf.ReadLimit(r, 64); err != nil {
				return err
			}
		default:
			p.fl.URI = append(p.fl.URI, p.buf.B[:i]...)
			p.buf.Discard(0, i+1)
			return nil
		}
	}
}

func (p *Probe) ParseVersion(r io.Reader) (version.Version, error) {
	for {
		i := p.buf.Index(char.CRLF) // HTTP/1.1\r\n
		switch i {
		case 0:
			return version.Version{}, httperr.ParseHTTPVersionErr
		case -1:
			size := 8 - p.buf.Len()
			if size <= 0 {
				return version.Version{}, httperr.RequestURITooLong
			}
			p.ExpandBuffer(size)
			if _, err := p.buf.ReadLimit(r, size); err != nil {
				return version.Version{}, err
			}
		default:
			ver, ok := version.Parse(p.buf.B[:i])
			if !ok {
				return version.Version{}, httperr.ParseHTTPVersionErr
			}

			p.fl.URI = append(p.fl.URI, p.buf.B[:i]...)
			p.buf.Discard(0, i+2)
			return ver, nil
		}
	}
}

func (p *Probe) ParseMethod(r io.Reader) error {
	if p.buf.Len() == 0 {
		return httperr.ParseHTTPMethodErr
	}

	for {
		i := p.buf.IndexByte(char.Space)

		switch i {
		case 0:
			return httperr.ParseHTTPMethodErr
		case -1:
			size := p.cfg.GetMaxRequestMethodSize() - p.buf.Len()
			if size <= 0 {
				return httperr.ParseHTTPMethodErr
			}
			p.ExpandBuffer(size)
			if _, err := p.buf.ReadLimit(r, size); err != nil {
				return err
			}
		default:
			p.fl.Method = method.Parse(p.buf.B[:i])
			p.buf.Discard(0, i+1)
			return nil
		}
	}
}

type FirstLineWriter struct {
}

var ConnectionPreface = []byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")

const ConnectionPrefaceLength = 24

func (f FirstLineWriter) Write(p []byte) (n int, err error) {
	// TODO implement me
	panic("implement me")
}

func (f FirstLineWriter) Close() error {
	// TODO implement me
	panic("implement me")
}
