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

	if ok {
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
	return r.Header.init()
}

func (s *Server) ParseRequestHeader(conn net.Conn, a *arena.Arena) (req *Request, err error) {
	req = arena.New[Request](a)

	var (
		b = arena.MakeSlice[byte](a, s.MaxRequestHeaderSize, s.MaxRequestHeaderSize)
		n int
	)

	if n, err = conn.Read(b); err != nil {
		return
	}
	buf := b[:n]

	if i := bytes.Index(buf, char.CRLF); i == -1 {
		return req, httperr.RequestURITooLong
	} else if err = req.parseFirstLine(buf[:i]); err == nil {
		buf = buf[i+2:]
	} else {
		return
	}

	ei := bytes.Index(buf, char.CRLF2) // end position index
	if ei == -1 {
		return req, httperr.RequestHeaderFieldsTooLarge
	}
	if err = req.parseHeaders(util.Split(a, buf[:ei], char.CRLF)); err != nil {
		return
	}

	// if len(req.Header.ContentLength) == 0 && req.Method {
	//
	// }

	// chunked
	if len(req.Header.TransferEncoding) != 0 && bytes.EqualFold(req.Header.TransferEncoding, char.Chunked) {

	}

	req.Header.TransferEncoding.Each(func(val []byte) bool {
		// 7\r\n
		// Mozilla\r\n
		// 11\r\n
		// Developer Network\r\n
		// 0\r\n
		// \r\n

		return true
	})

	if req.Method.IsTrace() { // TRACE request MUST NOT include an entity.
		_, _ = util.Copy(io.Discard, conn)
		return
	}

	// 	if req.MayContinue() {
	//		// 'Expect: 100-continue' header found. Let the caller deciding
	//		// whether to read request body or
	//		// to return StatusExpectationFailed.
	//		return nil
	//	}

	if length := req.Header.ContentLength.Length(); length == 0 {
		// Chunked
		req.Body = bodyio.Null()

	} else if length > s.MaxRequestBodySize {
		return nil, httperr.BodyTooLarge // body is too large
	} else if length > s.StreamRequestBodySize {
		req.Body, err = bodyio.File(a, buf[ei+4:], conn, length)
	} else if s.StreamRequestBodySize < 0 {
		req.Body, err = bodyio.Block(a, buf[ei+4:], conn, length)
	} else {
		req.Body, err = bodyio.Memory(a, buf[ei+4:], conn, length)
	}

	return
}
