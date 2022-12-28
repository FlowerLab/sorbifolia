//go:build goexperiment.arenas

package http

import (
	"arena"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"go.x2ox.com/sorbifolia/http/httperr"
)

type Handler func(ctx *Context)

const (
	defaultMaxRequestHeaderSize = 4 * 1024
	defaultMaxRequestBodySize   = 4 * 1024 * 1024
	defaultServerName           = "Sorbifolia"
	defaultUserAgent            = defaultServerName
)

var (
	defaultServerNameBytes = []byte(defaultServerName)
	// defaultReadTimeout     = time.Second *
)

type Server struct {
	Name []byte

	MaxRequestHeaderSize  int   // 最大允许的头大小，包括首行和 \r\n
	MaxRequestBodySize    int64 // 最大允许的 Body 大小
	StreamRequestBodySize int64 // 最大允许内存读入的 Body 大小

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

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

func (s *Server) serverName() []byte {
	if len(s.Name) != 0 {
		return s.Name
	}
	return defaultServerNameBytes
}

func (s *Server) handle(conn net.Conn) {
	go func() {

		// conn.SetReadDeadline(coarsetime.Now().Add(s.ReadTimeout))
		// conn.SetWriteDeadline(coarsetime.Now().Add(s.WriteTimeout))

		// defer conn.Close()
		a := arena.NewArena()
		req, err := s.ParseRequestHeader(conn, a)
		if err != nil {
			fmt.Println(err)
			conn.Write(req.ver.Bytes())

			switch err {
			case httperr.RequestHeaderFieldsTooLarge:
				conn.Write([]byte(" 431 Request Header Fields Too Large\r\n"))
			case httperr.BodyTooLarge:
				conn.Write([]byte(" 413 Request Entity Too Large\r\n"))
			default:
				conn.Write([]byte(" 500 Internal Server Error\r\n"))
			}

			conn.Write([]byte("Server: "))
			conn.Write(s.serverName())
			conn.Write([]byte("\r\nContent-Length: 0\r\n\r\n"))

			a.Free()
			return
		}

		ctx := s.newCtx(a, conn, req)
		atomic.AddInt64(&s.concurrency, 1)

		s.Handler(ctx)
		// send response
		_, _ = conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Server: "))
		conn.Write(s.serverName())
		conn.Write([]byte("\r\nContent-Length: 0\r\n\r\n"))

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
