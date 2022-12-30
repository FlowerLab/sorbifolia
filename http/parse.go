//go:build goexperiment.arenas

package http

import (
	"arena"
	"bytes"
	"io"
	"net"

	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bodyio"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

func (r *Request) parseFirstLine(b []byte) error {
	var (
		ok  bool
		arr = util.Split(r.a, b, char.Spaces)
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
	r.Header.KVs = arena.MakeSlice[KV](r.a, len(arr), len(arr))
	for i, v := range arr {
		r.Header.KVs[i].ParseHeader(v)
	}
	return r.Header.RawParse()
}

func (r *Request) Decode(s *Server, conn net.Conn) error {
	var (
		b      = arena.MakeSlice[byte](r.a, s.MaxRequestHeaderSize, s.MaxRequestHeaderSize)
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
	if err = r.parseHeaders(util.Split(r.a, buf[:ei], char.CRLF)); err != nil {
		return err
	}

	if r.Method.IsTrace() { // TRACE request MUST NOT include an entity.
		_, _ = util.Copy(io.Discard, conn)
		return nil
	}

	if bytes.Equal(r.Header.Get([]byte("Expect")).Val(), []byte("100-continue")) {
		r.Body = bodyio.Null()
	} else if length := r.Header.ContentLength.Length(); length == 0 {
		if bytes.Equal(r.Header.TransferEncoding, char.Chunked) {
			r.Body, err = bodyio.Chunked(r.a, buf[ei+4:], conn, int(s.MaxRequestBodySize))
		} else {
			r.Body = bodyio.Null()
		}
	} else if length > s.MaxRequestBodySize {
		err = httperr.BodyTooLarge // body is too large
	} else if length > s.StreamRequestBodySize {
		r.Body, err = bodyio.File(r.a, buf[ei+4:], conn, length)
	} else if s.StreamRequestBodySize < 0 {
		r.Body, err = bodyio.Block(r.a, buf[ei+4:], conn, length)
	} else {
		r.Body, err = bodyio.Memory(r.a, buf[ei+4:], conn, length)
	}

	return err
}
