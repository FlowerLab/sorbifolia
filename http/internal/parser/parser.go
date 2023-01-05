package parser

import (
	"bytes"
	"io"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
)

var http09 = []byte("HTTP/0.9")

type ChunkedTransfer func() (setTrailerHeader, setChunked func(b []byte) error)

type RequestParser struct {
	SetMethod  func([]byte) error
	SetURI     func([]byte) error
	SetVersion func([]byte) error
	SetHeaders func([]byte) (chunked ChunkedTransfer, length int64, err error)

	isChunked        bool
	setTrailerHeader func([]byte) error
	setChunked       func([]byte) error
	bodyLength       int64
	body             io.ReadCloser

	Limit httpconfig.Config
	state State
	err   error
	buf   bufpool.Buffer
}

type State uint8

const (
	ReadMethod State = iota
	ReadURI
	ReadVersion
	ReadHeader
	ReadBody
	ReadBodyChunked
	ReadBodyChunkedData
	ReadBodyChunkedEnd

	END
)

func (r *RequestParser) Write(p []byte) (n int, err error) {
	pLen := len(p)

	for len(p) > 0 {
		switch r.state {
		case ReadMethod:
			n, err = r.parseMethod(p)
		case ReadURI:
			n, err = r.parseURI(p)
		case ReadVersion:
			n, err = r.parseVersion(p)
		case ReadHeader:
			n, err = r.parseHeader(p)
		case ReadBody:
			n, err = r.parseBody(p)
		case ReadBodyChunked, ReadBodyChunkedData, ReadBodyChunkedEnd:
			n, err = r.parseBodyChunked(p)
		default:
			break
		}

		if err != nil {
			return 0, err
		}
		p = p[:n]
	}

	return pLen, err
}

func (r *RequestParser) parseBody(p []byte) (n int, err error) {
	// 看存放在什么地方
	return
}

func (r *RequestParser) parseBodyChunked(p []byte) (n int, err error) {
	var (
		i   = bytes.Index(p, char.CRLF) // Key: Value\r\n\r\nBody
		buf = &r.buf
	)

	if i == -1 || i == 1 { // buf[\r], p[\n\r\n] -> i == 1
		if length := buf.Len(); p[0] == '\n' && length > 0 && buf.B[length-1] == '\r' { // Check "\r\n" is straddles the buffer.
			i = 0  // The data in buf is enough, no need to read again
			n = -1 // Two bytes will be discarded later
			buf.B = buf.B[:length-1]
		}
	}

	switch i {
	case 0:
	case -1: // TODO add size limit
		return buf.Write(p)
	default:
		n, _ = buf.Write(p[:i])
	}

	switch r.state {
	case ReadBodyChunked:
		r.state++ // ReadBodyChunkedData. Has length, read data
		if buf.Len() == 1 && buf.B[0] == '0' {
			r.state++ // ReadBodyChunkedEnd. No length, read TrailerHeader or end
		}
		buf.Reset()
	case ReadBodyChunkedData:
		err = r.setChunked(buf.Bytes())
		r.state--
		buf.Reset()
	case ReadBodyChunkedEnd:
		if buf.Len() == 0 { // end
			r.state = END
			err = r.setTrailerHeader(buf.Bytes())
			buf.Reset()
		} else {
			_, _ = buf.Write(char.CRLF)
		}
	}

	n += 2 // Discard four bytes

	return
}

func (r *RequestParser) parseHeader(p []byte) (n int, err error) {
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
		if size := r.Limit.GetMaxRequestHeaderSize(); size > 0 && buf.Len()+len(p) > size {
			return 0, httperr.RequestHeaderFieldsTooLarge
		}
		return buf.Write(p)
	default:
		if size := r.Limit.GetMaxRequestHeaderSize(); size > 0 && buf.Len()+i > size {
			return 0, httperr.RequestHeaderFieldsTooLarge
		}
		if n, err = buf.Write(p[:i]); err != nil {
			return 0, err
		}
	}

	n += 4 // Discard four bytes
	r.state++

	var ct ChunkedTransfer
	if ct, r.bodyLength, err = r.SetHeaders(buf.Bytes()); err == nil && ct != nil {
		r.isChunked = true
		r.setTrailerHeader, r.setChunked = ct()
	}

	return
}

func (r *RequestParser) parseMethod(p []byte) (n int, err error) {
	var (
		i   = bytes.IndexByte(p, char.Space) // Method URI HTTP/1.1
		buf = &r.buf
	)

	switch i {
	case 0:
	case -1: // has not space
		if size := r.Limit.GetMaxRequestMethodSize(); size > 0 && buf.Len()+len(p) > size {
			return 0, httperr.ParseHTTPMethodErr
		}
		return buf.Write(p)
	default:
		if size := r.Limit.GetMaxRequestMethodSize(); size > 0 && buf.Len()+i > size {
			return 0, httperr.ParseHTTPMethodErr
		}
		if n, err = buf.Write(p[:i]); err != nil {
			return 0, err
		}
	}

	n++       // Discard a byte, it's a space
	r.state++ // Continue to read URI

	err = r.SetMethod(buf.Bytes())
	buf.Reset()

	return
}

func (r *RequestParser) parseURI(p []byte) (n int, err error) {
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
		if size := r.Limit.GetMaxRequestURISize(); size > 0 && buf.Len()+len(p) > size {
			return 0, httperr.RequestURITooLong
		}
		return buf.Write(p)
	default:
		if size := r.Limit.GetMaxRequestURISize(); size > 0 && buf.Len()+i > size {
			return 0, httperr.RequestURITooLong
		}
		if n, err = buf.Write(p[:i]); err != nil {
			return 0, err
		}
	}

	n++ // Discard a byte, it's a space
	r.state++

	err = r.SetURI(buf.Bytes())
	buf.Reset()

	if is09 {
		n++                      // Discard two bytes
		_ = r.SetVersion(http09) // this can't go wrong
		r.state = END            // HTTP/0.9 no header and body
	}

	return
}

func (r *RequestParser) parseVersion(p []byte) (n int, err error) {
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
		if size := r.Limit.GetMaxRequestURISize(); size > 0 && buf.Len()+len(p) > size {
			return 0, httperr.ParseHTTPMethodErr
		}
		return buf.Write(p)
	default:
		if size := r.Limit.GetMaxRequestURISize(); size > 0 && buf.Len()+i > size {
			return 0, httperr.ParseHTTPMethodErr
		}
		if n, err = buf.Write(p[:i]); err != nil {
			return 0, err
		}
	}

	n += 2                          // Discard two bytes
	err = r.SetVersion(buf.Bytes()) // this can't go wrong
	r.state++

	return
}

func (r *RequestParser) Close() error                     { panic("implement me") }
func (r *RequestParser) Read(p []byte) (n int, err error) { panic("implement me") }
func (r *RequestParser) Reset()                           {}
