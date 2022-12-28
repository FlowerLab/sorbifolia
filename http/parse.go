//go:build goexperiment.arenas

package http

import (
	"arena"
	"bytes"
	"fmt"
	"net"

	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bodyio"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

func (s *Server) ParseRequestHeader(conn net.Conn, a *arena.Arena) (req *Request, err error) {
	req = arena.New[Request](a)

	buf := arena.MakeSlice[byte](a, s.MaxRequestHeaderSize, s.MaxRequestHeaderSize)

	var n int
	if n, err = conn.Read(buf); err != nil {
		return
	}
	buf = buf[:n]

	idx := bytes.Index(buf, char.CRLF) // first line
	if idx == -1 {
		err = fmt.Errorf("parsing protocol header error")
		return
	}
	arr := util.Split(a, buf[:idx], char.Spaces)
	if len(arr) != 3 {
		err = fmt.Errorf("parsing protocol header error")
		return
	}

	req.Method = method.Parse(util.ToUpper(arr[0]))
	req.Header.RequestURI = arr[1]
	req.ver, _ = version.Parse(arr[2])

	ei := bytes.Index(buf, char.CRLF2) // end position index
	if ei == -1 {
		err = httperr.RequestHeaderFieldsTooLarge
		return
	}

	arr = util.Split(a, buf[idx+2:ei], char.CRLF) // header
	kvs := arena.MakeSlice[KV](a, len(arr), len(arr))
	for i, v := range arr {
		kvs[i].ParseHeader(v)
	}
	req.Header.KVs = (*KVs)(&kvs)

	if err = req.Header.init(); err != nil {
		return
	}
	// if len(req.Header.ContentLength) == 0 && req.Method {
	//
	// }

	// TODO length, check buf[ei+4:n] and conn
	if length := req.Header.ContentLength.Length(); length == 0 {
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
