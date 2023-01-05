package http

import (
	"io"
	"net"
	"sync/atomic"
	"time"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
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
	Config httpconfig.Config
	// Name   []byte

	// MaxRequestURISize     int   // 最大首行大小
	// MaxRequestHeaderSize  int   // 最大允许的头大小，包括首行和 \r\n
	// MaxRequestBodySize    int64 // 最大允许的 Body 大小
	// StreamRequestBodySize int64 // 最大允许内存读入的 Body 大小
	//
	// ReadTimeout  time.Duration
	// WriteTimeout time.Duration
	//
	// MaxIdleWorkerDuration              time.Duration
	// SleepWhenConcurrencyLimitsExceeded time.Duration

	Concurrency int
	Handler     Handler

	connCount   uint64
	concurrency uint32
	done        chan struct{}
}

func (s *Server) Listen() {}

func (s *Server) serveConnCleanup() {
	// atomic.AddInt32(&s.open, -1)
	atomic.AddUint32(&s.concurrency, ^uint32(0))
}

func (s *Server) fastWriteCode(w io.Writer, ver version.Version, code status.Status) error {
	if _, err := w.Write(ver.Bytes()); err != nil {
		return err
	}
	if _, err := w.Write(char.Spaces); err != nil {
		return err
	}
	if _, err := w.Write(code.Bytes()); err != nil {
		return err
	}
	if _, err := w.Write(char.CRLF); err != nil {
		return err
	}
	if _, err := w.Write([]byte("Connection: close\r\nServer: ")); err != nil {
		return err
	}
	if _, err := w.Write(s.Config.GetName()); err != nil {
		return err
	}
	_, err := w.Write([]byte("\r\nContent-Length: 0\r\n\r\n"))
	return err
}

func (s *Server) handle(conn net.Conn) error {
	defer s.serveConnCleanup()
	atomic.AddUint32(&s.concurrency, 1)

	// conn.SetReadDeadline(coarsetime.Now().Add(s.ReadTimeout))
	// conn.SetWriteDeadline(coarsetime.Now().Add(s.WriteTimeout))

	ctx := &Context{} // TODO: pool
	ctx.c = conn
	ctx.s = s

	ctx.id = atomic.AddUint64(&s.connCount, 1)
	ctx.time = time.Now()
	ctx.addr = conn.RemoteAddr()

	defer func() {
		ctx.c = nil
		_ = ctx.Request.Body.Close()
	}()

	ctx.Request.parse(conn)

	s.Handler(ctx)

	ctx.Response.Header.set(KV{
		K: char.Server,
		V: s.Config.GetName(),
	})
	ctx.Response.Header.set(KV{
		K: char.Date,
		V: util.GetDate(),
	})

	// var resp io.Reader
	resp, err := ctx.Response.Encode(ctx.Request.ver)
	if err != nil {
		_ = s.fastWriteCode(conn, ctx.Request.ver, status.InternalServerError)
		return nil
	}
	if _, err = util.Copy(conn, resp); err == io.EOF {
		err = nil
	}

	return err
}

func (s *Server) getConcurrency() int {
	if n := s.Concurrency; n > 0 {
		return n
	}
	return defaultConcurrency
}

func (s *Server) Serve(ln net.Listener) error {
	wp := &workerpool.WorkerPool{
		MaxWorkersCount:       s.getConcurrency(),
		MaxIdleWorkerDuration: s.Config.MaxIdleWorkerDuration,
	}
	wp.WorkerFunc = func(c net.Conn) error {
		err := s.handle(c)
		if err == nil {
			wp.SetConnState(c, workerpool.StateIdle)
		}
		return err
	}

	wp.Start()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		wp.SetConnState(conn, workerpool.StateNew)
		if !wp.Serve(conn) {
			_ = s.fastWriteCode(conn, version.Version{
				Major: 1, // read first line
				Minor: 1,
			}, status.ServiceUnavailable)
			_ = conn.Close()
			wp.SetConnState(conn, workerpool.StateClosed)

			if s.Config.SleepWhenConcurrencyLimitsExceeded > 0 {
				time.Sleep(s.Config.SleepWhenConcurrencyLimitsExceeded)
			}
		}

		conn = nil
	}
}
