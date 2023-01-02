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

func parseBody(r *Request, s *Server, pr []byte, read io.Reader, max int) (err error) {
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
			r.Body, err = bodyio.Chunked(pr, read, max)
		} else {
			r.Body = bodyio.Null()
		}
		return
	}

	if length > int64(max) {
		return httperr.BodyTooLarge // body is too large
	} else if length > s.StreamRequestBodySize {
		r.Body, err = bodyio.File(pr, read, length)
	} else if s.StreamRequestBodySize < 0 {
		r.Body, err = bodyio.Block(pr, read, length)
	} else {
		r.Body, err = bodyio.Memory(pr, read, length)
	}

	return err
}

func parseHeaders(r *Request, pr []byte, read io.Reader, max int) error {
	buf := bufpool.Acquire(max)
	defer bufpool.Release(buf)

	wn, _ := buf.Write(pr)

	_, err := util.Copy(buf, io.LimitReader(read, int64(max-wn)))
	if err != nil && err != io.EOF {
		return err
	}

	b := buf.Bytes()

	idx := bytes.Index(b, char.CRLF2)
	if idx < 0 {
		return httperr.RequestHeaderFieldsTooLarge
	}

	body := b[idx+4:]
	b = b[:idx+2]
	_ = body

	r.Header.KVs.preAlloc(bytes.Count(b, char.CRLF))

	for idx = bytes.Index(b, char.CRLF); len(b) > 0; {
		r.Header.KVs.addHeader(b[:idx])
		b = b[idx+2:]
	}

	return r.Header.RawParse()
}

func parseFirstLine(r *Request, read io.Reader, max int) error {
	buf := bufpool.Acquire(max)
	defer bufpool.Release(buf)

	n, err := util.Copy(buf, io.LimitReader(read, int64(max)))
	if err != nil && err != io.EOF {
		return err
	}

	idx := bytes.Index(buf.B[:n], char.CRLF)
	if idx < 0 {
		return httperr.RequestURITooLong
	}

	b := buf.B[:idx]

	h := buf.B[idx+2 : n]
	hb := make([]byte, len(h))
	copy(hb, h)

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
		b      = make([]byte, s.MaxRequestHeaderSize)
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
			r.Body, err = bodyio.Chunked(buf[ei+4:], conn, int(s.MaxRequestBodySize))
		} else {
			r.Body = bodyio.Null()
		}
	} else if length > s.MaxRequestBodySize {
		err = httperr.BodyTooLarge // body is too large
	} else if length > s.StreamRequestBodySize {
		r.Body, err = bodyio.File(buf[ei+4:], conn, length)
	} else if s.StreamRequestBodySize < 0 {
		r.Body, err = bodyio.Block(buf[ei+4:], conn, length)
	} else {
		r.Body, err = bodyio.Memory(buf[ei+4:], conn, length)
	}

	return err
}
