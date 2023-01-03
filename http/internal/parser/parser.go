package parser

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
)

var http09 = []byte("HTTP/0.9")

type RequestParser struct {
	SetMethod  func([]byte) error
	SetURI     func([]byte) error
	SetVersion func([]byte) error
	SetHeaders func([]byte) error
	SetBody    func([]byte) error

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
		case ReadBody:
			n, err = r.parseMethod(p)
		}

		if err != nil {
			return 0, err
		}
		p = p[:n]
	}

	return pLen, err
}

func (r *RequestParser) parseHeader(p []byte) (n int, err error) {

	var (
		i   = bytes.Index(p, char.CRLF2) // Key: Value\r\n\r\nBody
		buf = &r.buf
	)
	if i == -1 {
		// Check "\r\n\r\n" is straddles the buffer.
		var _b [6]byte

		length := buf.Len()
		if length >= 3 {
			n = copy(_b[:], buf.B[length-3:])
		} else {
			n = copy(_b[:], buf.B)
		}
		copy(_b[n:], p)

		if i = bytes.Index(_b[:], char.CRLF2); i != -1 {
			n = i - n
			i = 0
			buf.B = buf.B[:n+1]
		}
	}

	if i == -1 { // has not CRLF2
		if size := r.Limit.GetMaxRequestURISize(); size > 0 && buf.Len()+len(p) > size {
			// check the len(p) to avoid exceeding the limit
			return 0, httperr.RequestHeaderFieldsTooLarge
		}
		return buf.Write(p)
	}

	if i > 0 { // There are also bytes in p
		n, err = buf.Write(p[:i])
	}
	n += 4 // Discard four bytes
	err = r.SetHeaders(buf.Bytes())
	r.state++

	return
}

func (r *RequestParser) parseMethod(p []byte) (n int, err error) {
	var (
		i   = bytes.IndexByte(p, char.Space) // Method URI HTTP/1.1
		buf = r.buf
	)

	if i == -1 { // has not space
		if size := r.Limit.GetMaxRequestMethodSize(); size > 0 && buf.Len()+len(p) > size {
			// check the len(p) to avoid exceeding the limit
			return 0, httperr.ParseHTTPMethodErr
		}
		return buf.Write(p)
	}

	if i > 0 { // There are also bytes in p
		n, err = buf.Write(p[:i])
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
		buf  = r.buf
		is09 = false
	)
	if i == -1 {
		if i = bytes.Index(p, char.CRLF); i >= 0 { // HTTP/0.9 no version
			is09 = true
		} else if p[0] == '\n' && buf.Len() > 0 && buf.B[buf.Len()-1] == 'r' { // Check "\r\n" is straddles the buffer.
			i = 0  // The data in buf is enough, no need to read again
			n = -1 // Two bytes will be discarded later
		}
	}

	if i == -1 { // has not space
		if size := r.Limit.GetMaxRequestURISize(); size > 0 && buf.Len()+len(p) > size {
			// check the len(p) to avoid exceeding the limit
			return 0, httperr.RequestURITooLong
		}
		return buf.Write(p)
	}

	if i > 0 { // There are also bytes in p
		n, err = buf.Write(p[:i])
	}
	if is09 {
		n += 2                   // Discard two bytes
		_ = r.SetVersion(http09) // this can't go wrong
		r.state = END            // HTTP/0.9 no header and body
	} else {
		n++ // Discard a byte, it's a space
		r.state++
	}

	err = r.SetMethod(buf.Bytes())
	buf.Reset()

	return
}

func (r *RequestParser) parseVersion(p []byte) (n int, err error) {
	var (
		i   = bytes.Index(p, char.CRLF) // Method URI HTTP/1.1
		buf = r.buf
	)
	if i == -1 {
		if p[0] == '\n' && buf.Len() > 0 && buf.B[buf.Len()-1] == 'r' { // Check "\r\n" is straddles the buffer.
			i = 0  // The data in buf is enough, no need to read again
			n = -1 // Two bytes will be discarded later
		}
	}

	if i == -1 { // has not CRLF
		if size := r.Limit.GetMaxRequestURISize(); size > 0 && buf.Len()+len(p) > size {
			// check the len(p) to avoid exceeding the limit
			return 0, httperr.ParseHTTPVersionErr
		}
		return buf.Write(p)
	}

	if i > 0 { // There are also bytes in p
		n, err = buf.Write(p[:i])
	}
	n += 2                          // Discard two bytes
	err = r.SetVersion(buf.Bytes()) // this can't go wrong
	r.state++

	return
}

func (r *RequestParser) Close() error                     { panic("implement me") }
func (r *RequestParser) Read(p []byte) (n int, err error) { panic("implement me") }
func (r *RequestParser) Reset()                           {}
