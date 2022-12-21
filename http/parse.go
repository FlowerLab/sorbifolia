//go:build goexperiment.arenas

package http

import (
	"arena"
	"bytes"
	"fmt"
	"net"

	"go.x2ox.com/sorbifolia/http/internal/bodyio"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

func (s *Server) ParseRequestHeader(conn net.Conn, a *arena.Arena) (*Request, error) {
	buf := arena.MakeSlice[byte](a, s.MaxRequestHeaderSize, s.MaxRequestHeaderSize)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	buf = buf[:n]

	idx := bytes.Index(buf, char.CRLF) // first line
	if idx == -1 {
		return nil, fmt.Errorf("parsing protocol header error")
	}
	arr := util.Split(a, buf[:idx], char.Spaces)
	if len(arr) != 3 {
		return nil, fmt.Errorf("parsing protocol header error")
	}

	req := arena.New[Request](a)
	req.Method = method.Parse(util.ToUpper(arr[0]))
	req.Header.RequestURI = arr[1]
	req.ver, _ = version.Parse(arr[2])

	ei := bytes.Index(buf, char.CRLF2) // end position index
	if ei == -1 {
		return nil, fmt.Errorf("413 Entity Too Large")
	}

	arr = util.Split(a, buf[idx+2:ei], char.CRLF) // header
	kvs := arena.MakeSlice[KV](a, len(arr), len(arr))
	for i, v := range arr {
		kvs[i].ParseHeader(v)
	}
	req.Header.KVs = (*KVs)(&kvs)

	// TODO length, check buf[ei+4:n] and conn
	if req.Header.ContentLength == 0 {
		req.Body = bodyio.Null()
	} else if req.Header.ContentLength > s.MaxRequestBodySize {
		return nil, err // body is too large
	} else if req.Header.ContentLength > s.StreamRequestBodySize {
		req.Body, err = bodyio.File(a, buf[ei+4:], conn, req.Header.ContentLength)
	} else if s.StreamRequestBodySize < 0 {
		req.Body, err = bodyio.Block(a, buf[ei+4:], conn, req.Header.ContentLength)
	} else {
		req.Body, err = bodyio.Memory(a, buf[ei+4:], conn, req.Header.ContentLength)
	}
	if err != nil {
		return nil, err
	}

	return req, nil
}
