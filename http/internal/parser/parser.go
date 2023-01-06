package parser

import (
	"bytes"
	"errors"
	"io"
	"sync"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
)

var (
	http09             = []byte("HTTP/0.9")
	_RequestParserPool = sync.Pool{}
)

func AcquireRequestParserAAAA() *RequestParser {
	if v := _RequestParserPool.Get(); v != nil {
		return v.(*RequestParser)
	}
	return &RequestParser{}
}

func AcquireRequestParser(
	setMethod, setURI, setVersion func([]byte) error,
	setHeader func([]byte) (length int, err error),
) *RequestParser {
	if v := _RequestParserPool.Get(); v != nil {
		br := v.(*RequestParser)
		br.SetMethod = setMethod
		br.SetURI = setURI
		br.SetVersion = setVersion
		br.SetHeaders = setHeader
		return br
	}
	return &RequestParser{
		SetMethod:  setMethod,
		SetURI:     setURI,
		SetVersion: setVersion,
		SetHeaders: setHeader,
	}
}

func ReleaseRequestParser(r *RequestParser) { r.Reset(); _RequestParserPool.Put(r) }

type ChunkedTransfer func() (setTrailerHeader, setChunked func(b []byte) error)

type RequestParser struct {
	Limit      httpconfig.Config
	SetMethod  func([]byte) error
	SetURI     func([]byte) error
	SetVersion func([]byte) error
	SetHeaders func([]byte) (length int, err error)
	BW         io.WriteCloser

	setTrailerHeader, setChunked func([]byte) error
	bodyLength                   int

	state State
	buf   bufpool.Buffer
	rp    int
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
		case ReadHeader:
			n, err = r.parseHeader(p)
		case ReadBody:
			n, err = r.parseBody(p)
		default:
			break
		}

		if err != nil {
			return 0, err
		}
		p = p[n:]
	}

	if r.state == END {
		err = io.EOF
	}

	return pLen, err
}

func (r *RequestParser) parseBody(p []byte) (n int, err error) {
	if n, err = r.BW.Write(p); err != nil {
		if r.bodyLength < 0 { // chunked mode
			if err == io.EOF {
				if cErr := r.BW.Close(); cErr != nil {
					err = cErr
				}
				r.state = END
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
		err = io.EOF
		r.state = END
	}
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
		n, _ = buf.Write(p[:i])
	}

	n += 4 // Discard four bytes
	r.state++

	if r.bodyLength, err = r.SetHeaders(buf.Bytes()); err != nil {
		return
	}
	if r.bodyLength == 0 { // Not has Body
		r.state = END
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
		n, _ = buf.Write(p[:i])
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
		n, _ = buf.Write(p[:i])
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
		n, _ = buf.Write(p[:i])
	}

	n += 2                          // Discard two bytes
	err = r.SetVersion(buf.Bytes()) // this can't go wrong
	r.state++
	buf.Reset()

	return
}

func (r *RequestParser) Close() error {
	r.buf.Reset()
	return nil
}

func (r *RequestParser) Read(p []byte) (n int, err error) {
	if r.state != END {
		return 0, errors.New("? TODO: here we need to think about how to deal with")
	}
	if length := r.buf.Len(); length == 0 || length == r.rp {
		return 0, io.EOF
	}
	n = copy(p, r.buf.B[r.rp:])
	r.rp += n
	return
}

func (r *RequestParser) Reset() {
	r.SetMethod = nil
	r.SetURI = nil
	r.SetVersion = nil
	r.SetHeaders = nil
	r.setTrailerHeader = nil
	r.setChunked = nil
	r.state = ReadMethod
	r.buf.Reset()
	r.rp = 0
}
