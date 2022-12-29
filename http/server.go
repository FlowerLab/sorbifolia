//go:build goexperiment.arenas

package http

import (
	"arena"
	"io"
	"net"
	"sync/atomic"
	"time"

	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/workerpool"
	"go.x2ox.com/sorbifolia/http/status"
	"go.x2ox.com/sorbifolia/http/version"
)

type Handler func(ctx *Context)

const (
	defaultMaxRequestHeaderSize = 4 * 1024
	defaultMaxRequestBodySize   = 4 * 1024 * 1024
	defaultServerName           = "Sorbifolia"
	defaultUserAgent            = defaultServerName

	defaultConcurrency = 256 * 1024
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

	Concurrency                        int
	MaxIdleWorkerDuration              time.Duration
	SleepWhenConcurrencyLimitsExceeded time.Duration
	Handler                            Handler

	connCount   uint64
	concurrency uint32
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

func (s *Server) getServerName() []byte {
	if len(s.Name) != 0 {
		return s.Name
	}
	return defaultServerNameBytes
}

func (s *Server) serveConnCleanup() {
	// atomic.AddInt32(&s.open, -1)
	atomic.AddUint32(&s.concurrency, ^uint32(0))
}

func (s *Server) fastWriteCode(w io.Writer, ver version.Version, code status.Status) {
	_, _ = w.Write(ver.Bytes())
	_, _ = w.Write(char.Spaces)
	_, _ = w.Write(code.Bytes())
	_, _ = w.Write(char.CRLF)
	_, _ = w.Write([]byte("Connection: close\r\nServer: "))
	_, _ = w.Write(s.getServerName())
	_, _ = w.Write([]byte("\r\nContent-Length: 0\r\n\r\n"))
}

func (s *Server) handle(conn net.Conn) error {
	defer s.serveConnCleanup()
	atomic.AddUint32(&s.concurrency, 1)

	// conn.SetReadDeadline(coarsetime.Now().Add(s.ReadTimeout))
	// conn.SetWriteDeadline(coarsetime.Now().Add(s.WriteTimeout))

	a := arena.NewArena()
	req, err := s.ParseRequestHeader(conn, a)
	if err != nil {
		switch err {
		case httperr.RequestHeaderFieldsTooLarge:
			s.fastWriteCode(conn, req.ver, status.RequestEntityTooLarge)
		case httperr.BodyTooLarge:
			s.fastWriteCode(conn, req.ver, status.RequestEntityTooLarge)
		default:
			s.fastWriteCode(conn, req.ver, status.InternalServerError)
		}

		a.Free()
		return nil
	}

	ctx := s.newCtx(a, conn, req)

	s.Handler(ctx)
	s.fastWriteCode(conn, req.ver, status.OK)

	ctx.c = nil
	if ctx.robbery {
		return nil
	}
	a.Free()

	return nil
}

func (s *Server) getConcurrency() int {
	n := s.Concurrency
	if n <= 0 {
		n = defaultConcurrency
	}
	return n
}

func (s *Server) Serve(ln net.Listener) error {
	wp := &workerpool.WorkerPool{
		WorkerFunc:            s.handle,
		MaxWorkersCount:       s.getConcurrency(),
		MaxIdleWorkerDuration: s.MaxIdleWorkerDuration,
	}
	wp.Start()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		wp.SetConnState(conn, workerpool.StateNew)
		if !wp.Serve(conn) {
			// s.writeFastError(c, StatusServiceUnavailable,
			// 	"The connection cannot be served because Server.Concurrency limit exceeded")
			_ = conn.Close()
			wp.SetConnState(conn, workerpool.StateClosed)

			if s.SleepWhenConcurrencyLimitsExceeded > 0 {
				time.Sleep(s.SleepWhenConcurrencyLimitsExceeded)
			}
		}

		conn = nil
	}
}
