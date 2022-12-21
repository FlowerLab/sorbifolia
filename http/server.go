//go:build goexperiment.arenas

package http

import (
	"arena"
	"encoding/json"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type Handler func(ctx *Context)

const (
	defaultMaxRequestHeaderSize = 4 * 1024
	defaultMaxRequestBodySize   = 4 * 1024 * 1024
	defaultServerName           = "Sorbifolia"
	defaultUserAgent            = defaultServerName
)

type Server struct {
	Name []byte

	MaxRequestHeaderSize  int   // 最大允许的头大小，包括首行和 \r\n
	MaxRequestBodySize    int64 // 最大允许的 Body 大小
	StreamRequestBodySize int64 // 最大允许内存读入的 Body 大小

	Handler Handler

	connCount   uint64
	concurrency int64
	done        chan struct{}
}

func (s *Server) Listen() {}
func (s *Server) newCtx(a *arena.Arena, conn net.Conn, req *Request) *Context {
	ctx := arena.New[Context](a)
	ctx.a = a
	ctx.c = conn
	ctx.s = s

	ctx.id = atomic.AddUint64(&s.connCount, 1)
	ctx.time = time.Now()
	ctx.addr = conn.RemoteAddr()

	ctx.Request = req
	return ctx
}

func (s *Server) handle(conn net.Conn) {
	go func() {
		a := arena.NewArena()
		req, err := s.ParseRequestHeader(conn, a)
		if err != nil {
			fmt.Println(err)
			return
		}
		bbb, _ := json.Marshal(req)
		fmt.Println(string(bbb))

		ctx := s.newCtx(a, conn, req)
		atomic.AddInt64(&s.concurrency, 1)

		s.Handler(ctx)
		// send response
		_, _ = conn.Write([]byte("HTTP/1.1 200 OK\r\nserver: com\r\ncontent-length: 0\r\n\r\n"))

		atomic.AddInt64(&s.concurrency, -1)
		ctx.c = nil
		if ctx.robbery {
			return
		}
		a.Free()
	}()
}

func (s *Server) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		s.handle(conn)
	}
}
