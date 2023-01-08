package httpmessage

import (
	"bytes"
	"io"

	"go.x2ox.com/sorbifolia/http/httpbody"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

var http09 = []byte("HTTP/0.9")

func (r *Request) Write(p []byte) (n int, err error) {
	if !r.state.Writable() {
		return 0, io.EOF
	}
	pLen := len(p)

	for len(p) > 0 {
		switch r.state.Operate() {
		case _Init:
			r.state.SetOperate(_Method)
			continue
		case _Method:
			n, err = r.parseMethod(p)
		case _URI:
			n, err = r.parseURI(p)
		case _Version:
			n, err = r.parseVersion(p)
		case _Header:
			n, err = r.parseHeader(p)
		case _Body:
			n, err = r.parseBody(p)
		default:
			break
		}

		if err != nil {
			return 0, err
		}
		p = p[n:]
	}

	if r.state.IsClose() {
		err = io.EOF
	}

	return pLen, err
}

func (r *Request) parseBody(p []byte) (n int, err error) {
	if n, err = r.Body.Write(p); err != nil {
		if r.bodyLength < 0 { // chunked mode
			if err == io.EOF {
				if cErr := r.Body.Close(); cErr != nil {
					err = cErr
				}
				r.state.Close()
			}
		}
		return
	}
	if r.bodyLength < 0 { // chunked mode
		return
	}

	r.bodyLength -= n
	if r.bodyLength < 0 {
		err = io.ErrUnexpectedEOF // TODO err
	} else if r.bodyLength == 0 {
		if err = r.Body.Close(); err == nil {
			err = io.EOF
		}
		r.state.Close()
	}
	return
}

func (r *Request) parseHeader(p []byte) (n int, err error) {
	var (
		i   = bytes.Index(p, char.CRLF2) // \r\n
		buf = &r.buf
	)

	if length := buf.Len(); i == -1 && buf.Len()+len(p) >= 4 { // Check "\r\n\r\n" is straddles the buffer.
		var (
			_b         [6]byte
			copyLength = 3
		)

		if length >= copyLength {
			n = copy(_b[:], buf.B[length-copyLength:])
		} else {
			copyLength = length
			n = copy(_b[:], buf.B)
		}
		copy(_b[n:], p)

		if i = bytes.Index(_b[:], char.CRLF2); i != -1 {
			if idx := copyLength - i; idx > 0 {
				buf.B = buf.B[:length-idx]
			}
			n = i - n
			i = 0
		}
	}

	switch i {
	case 0:
	case -1:
		if size := r.cfg.GetMaxRequestHeaderSize(); size > 0 && buf.Len()+len(p) > size {
			return 0, httperr.RequestHeaderFieldsTooLarge
		}
		return buf.Write(p)
	default:
		if size := r.cfg.GetMaxRequestHeaderSize(); size > 0 && buf.Len()+i > size {
			return 0, httperr.RequestHeaderFieldsTooLarge
		}
		n, _ = buf.Write(p[:i])
	}

	n += 4 // Discard four bytes
	r.state.SetOperate(_Body)

	if err = r.Header.Parse(buf.Bytes()); err != nil {
		return
	}

	r.bodyLength = int(r.Header.ContentLength.Length())
	switch r.bodyLength {
	case 0:
		r.state.Close()
		r.Body = httpbody.Null()
	case -1:
		c := httpbody.AcquireChunked()
		c.Data = make(chan []byte, 1)
		c.Header = make(chan []byte, 1)
		r.Body = c
	default:
		r.Body = httpbody.AcquireMemory()
	}

	return
}

func (r *Request) parseMethod(p []byte) (n int, err error) {
	var (
		i   = bytes.IndexByte(p, char.Space) // Method URI HTTP/1.1
		buf = &r.buf
	)

	switch i {
	case 0:
	case -1: // has not space
		if size := r.cfg.GetMaxRequestMethodSize(); size > 0 && buf.Len()+len(p) > size {
			return 0, httperr.ParseHTTPMethodErr
		}
		return buf.Write(p)
	default:
		if size := r.cfg.GetMaxRequestMethodSize(); size > 0 && buf.Len()+i > size {
			return 0, httperr.ParseHTTPMethodErr
		}
		n, _ = buf.Write(p[:i])
	}

	n++                      // Discard a byte, it's a space
	r.state.SetOperate(_URI) // Continue to read URI

	r.Method = method.Parse(buf.Bytes())
	buf.Reset()

	return
}

func (r *Request) parseURI(p []byte) (n int, err error) {
	var (
		i    = bytes.IndexByte(p, char.Space) // Method URI HTTP/1.1
		buf  = &r.buf
		is09 = false
	)
	if i == -1 {
		if i = bytes.Index(p, char.CRLF); i >= 0 { // HTTP/0.9 no version
			is09 = true
		} else if length := buf.Len(); p[0] == '\n' && length > 0 && buf.B[length-1] == '\r' { // Check "\r\n" is straddles the buffer.
			i = 0  // The data in buf is enough, no need to read again
			n = -1 // Two bytes will be discarded later
			buf.B = buf.B[:length-1]
		}
	}

	switch i {
	case 0:
	case -1: // has not space
		if size := r.cfg.GetMaxRequestURISize(); size > 0 && buf.Len()+len(p) > size {
			return 0, httperr.RequestURITooLong
		}
		return buf.Write(p)
	default:
		if size := r.cfg.GetMaxRequestURISize(); size > 0 && buf.Len()+i > size {
			return 0, httperr.RequestURITooLong
		}
		n, _ = buf.Write(p[:i])
	}

	n++ // Discard a byte, it's a space
	r.state.SetOperate(_Version)

	r.Header.RequestURI = append(r.Header.RequestURI, buf.Bytes()...)

	buf.Reset()

	if is09 {
		n++                                  // Discard two bytes
		r.Version, _ = version.Parse(http09) // this can't go wrong
		r.state.Close()                      // HTTP/0.9 no header and body
	}

	return
}

func (r *Request) parseVersion(p []byte) (n int, err error) {
	var (
		i   = bytes.Index(p, char.CRLF) // Method URI HTTP/1.1
		buf = &r.buf
	)
	if i == -1 {
		if p[0] == '\n' && buf.Len() > 0 && buf.B[buf.Len()-1] == 'r' { // Check "\r\n" is straddles the buffer.
			i = 0  // The data in buf is enough, no need to read again
			n = -1 // Two bytes will be discarded later
		}
	}

	switch i {
	case 0:
	case -1: // has not space
		if size := r.cfg.GetMaxRequestURISize(); size > 0 && buf.Len()+len(p) > size {
			return 0, httperr.ParseHTTPMethodErr
		}
		return buf.Write(p)
	default:
		if size := r.cfg.GetMaxRequestURISize(); size > 0 && buf.Len()+i > size {
			return 0, httperr.ParseHTTPMethodErr
		}
		n, _ = buf.Write(p[:i])
	}

	var ok bool
	if r.Version, ok = version.Parse(buf.Bytes()); !ok {
		err = httperr.ParseHTTPVersionErr
	}

	n += 2 // Discard two bytes
	r.state.SetOperate(_Header)
	buf.Reset()

	return
}
