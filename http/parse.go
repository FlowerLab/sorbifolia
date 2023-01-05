package http

import (
	"bytes"
	"io"
	"net"

	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bodyio"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

func (r *Request) parseFirstLine(b []byte) error {
	var (
		ok  bool
		arr = bytes.Split(b, char.Spaces)
	)

	switch len(arr) {
	case 2:
		r.Method = method.Parse(util.ToUpper(arr[0]))
		r.Header.RequestURI = arr[1]
		r.ver.Major, r.ver.Minor = 0, 9
		ok = true
	case 3:
		r.Method = method.Parse(util.ToUpper(arr[0]))
		r.Header.RequestURI = arr[1]
		r.ver, ok = version.Parse(arr[2])
	}

	if !ok {
		return httperr.ParseHTTPVersionErr
	}
	return nil
}

func (r *Request) parseHeaders(arr [][]byte) error {
	if len(arr) == 0 {
		return nil
	}
	r.Header.KVs = make([]KV, len(arr))
	for i, v := range arr {
		r.Header.KVs[i].ParseHeader(v)
	}
	return r.Header.RawParse()
}

func (r *Request) parseBody(s *Server, read io.Reader, buf *bufpool.Buffer, max int) (err error) {
	if r.Method.IsTrace() { // TRACE request MUST NOT include an entity.
		return nil // util.Copy(io.Discard, conn)
	}

	if bytes.Equal(r.Header.Get([]byte("Expect")).V, []byte("100-continue")) {
		r.Body = bodyio.Null()
		return
	}

	length := r.Header.ContentLength.Length()
	if length == 0 {
		if bytes.Equal(r.Header.TransferEncoding, char.Chunked) {
			r.Body, err = bodyio.Chunked(buf.Bytes(), read, max)
		} else {
			r.Body = bodyio.Null()
		}
		return
	}

	if length > int64(max) {
		return httperr.BodyTooLarge // body is too large
	} else if length > int64(s.Config.StreamRequestBodySize) {
		r.Body, err = bodyio.File(buf.Bytes(), read, length)
	} else if s.Config.StreamRequestBodySize < 0 {
		r.Body, err = bodyio.Block(buf.Bytes(), read, length)
	} else {
		r.Body, err = bodyio.Memory(buf.Bytes(), read, length)
	}

	return err
}

func (r *Request) _parseHeaders(read io.Reader, buf, body *bufpool.Buffer) error {
	if _, err := util.Copy(buf, read); err != nil && err != io.EOF {
		return err
	}

	b := buf.Bytes()
	idx := bytes.Index(b, char.CRLF2)
	if idx < 0 {
		return httperr.RequestHeaderFieldsTooLarge
	}

	body.B = append(body.B, buf.B[idx+4:]...)
	b = b[:idx+2]

	r.Header.KVs.preAlloc(bytes.Count(b, char.CRLF))

	for idx = bytes.Index(b, char.CRLF); len(b) > 0; {
		r.Header.KVs.addHeader(b[:idx])
		b = b[idx+2:]
	}

	return r.Header.RawParse()
}

func (r *Request) _parseFirstLine(read io.Reader, buf, header *bufpool.Buffer) error {
	if _, err := util.Copy(buf, read); err != nil && err != io.EOF {
		return err
	}

	b := buf.Bytes()
	idx := bytes.Index(buf.B, char.CRLF)
	if idx < 0 {
		return httperr.RequestURITooLong
	}

	header.B = append(header.B, buf.B[idx+2:]...)
	b = buf.B[:idx]

	if idx = bytes.IndexByte(b, char.Space); idx < 0 {
		return httperr.ParseHTTPMethodErr
	}
	r.Method = method.Parse(util.ToUpper(b[:idx]))
	b = b[idx+1:]

	if idx = bytes.IndexByte(b, char.Space); idx < 0 { // HTTP/0.9 no protocol header
		r.Header.RequestURI = append(r.Header.RequestURI, b...)
		r.ver.Major, r.ver.Minor = 0, 9
	} else {
		r.Header.RequestURI = append(r.Header.RequestURI, b[:idx]...)

		var ok bool
		if r.ver, ok = version.Parse(b[idx+1:]); !ok {
			return httperr.ParseHTTPVersionErr
		}
	}

	return nil
}

func (r *Request) Decode(s *Server, conn net.Conn) error {
	var (
		b      = make([]byte, s.Config.MaxRequestHeaderSize)
		n, err = conn.Read(b)
	)

	if err != nil {
		return err
	}
	buf := b[:n]

	if i := bytes.Index(buf, char.CRLF); i == -1 {
		return httperr.RequestURITooLong
	} else if err = r.parseFirstLine(buf[:i]); err != nil {
		return err
	} else {
		buf = buf[i+2:]
	}

	ei := bytes.Index(buf, char.CRLF2) // end position index
	if ei == -1 {
		return httperr.RequestHeaderFieldsTooLarge
	}
	if err = r.parseHeaders(bytes.Split(buf[:ei], char.CRLF)); err != nil {
		return err
	}

	if r.Method.IsTrace() { // TRACE request MUST NOT include an entity.
		_, _ = util.Copy(io.Discard, conn)
		return nil
	}

	if bytes.Equal(r.Header.Get([]byte("Expect")).V, []byte("100-continue")) {
		r.Body = bodyio.Null()
	} else if length := r.Header.ContentLength.Length(); length == 0 {
		if bytes.Equal(r.Header.TransferEncoding, char.Chunked) {
			r.Body, err = bodyio.Chunked(buf[ei+4:], conn, int(s.Config.MaxRequestBodySize))
		} else {
			r.Body = bodyio.Null()
		}
	} else if length > int64(s.Config.MaxRequestBodySize) {
		err = httperr.BodyTooLarge // body is too large
	} else if length > int64(s.Config.StreamRequestBodySize) {
		r.Body, err = bodyio.File(buf[ei+4:], conn, length)
	} else if s.Config.StreamRequestBodySize < 0 {
		r.Body, err = bodyio.Block(buf[ei+4:], conn, length)
	} else {
		r.Body, err = bodyio.Memory(buf[ei+4:], conn, length)
	}

	return err
}
