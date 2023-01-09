package http

import (
	"bytes"
	"io"
	"net"
	"sync/atomic"
	"time"

	"go.x2ox.com/sorbifolia/coarsetime"
	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/http/internal/workerpool"
	"go.x2ox.com/sorbifolia/http/status"
	"go.x2ox.com/sorbifolia/http/version"
)

type Handler func(ctx *Context)

type Server struct {
	Config httpconfig.Config

	Handler Handler

	connCount   uint64
	concurrency uint32
	done        chan struct{}
}

func (s *Server) Listen() {}

func (s *Server) serveConnCleanup() {
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

func (s *Server) getCtx(conn net.Conn) *Context {
	ctx := AcquireContext()
	ctx.c = conn
	ctx.s = s
	ctx.id = atomic.AddUint64(&s.connCount, 1)
	ctx.time = time.Now()
	ctx.addr = conn.RemoteAddr()

	return ctx
}

func (s *Server) serveConn(conn net.Conn) error {
	atomic.AddUint32(&s.concurrency, 1)
	defer s.serveConnCleanup()

	ctx := s.getCtx(conn)
	defer ReleaseContext(ctx)

	var (
		err error
	)

	for {
		if err = conn.SetReadDeadline(coarsetime.Now().Add(s.Config.GetReadTimeout())); err != nil {
			break
		}
		if _, err = util.Copy(&ctx.Request, conn); err != nil && err != io.EOF {
			break
		}
		s.Handler(ctx)

		if ctx.Response.Version.Null() {
			ctx.Response.Version = ctx.Request.Version
		}

		ctx.Response.Header.Set(char.Server, s.Config.GetName())
		ctx.Response.Header.Set(char.Date, util.GetDate())

		var buf = &bytes.Buffer{}

		if _, err = util.Copy(buf, &ctx.Response); err != nil && err != io.EOF {
			break
		}

		if _, err = util.Copy(conn, buf); err != nil && err != io.EOF {
			break
		}

		ctx.cleanup()
	}

	return err
}

func (s *Server) Serve(ln net.Listener) error {
	wp := &workerpool.WorkerPool{
		MaxWorkersCount:       s.Config.GetConcurrency(),
		MaxIdleWorkerDuration: s.Config.MaxIdleWorkerDuration,
	}
	wp.WorkerFunc = func(c net.Conn) error {
		err := s.serveConn(c)
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
			_ = conn.Close()
			wp.SetConnState(conn, workerpool.StateClosed)

			if s.Config.SleepWhenConcurrencyLimitsExceeded > 0 {
				time.Sleep(s.Config.SleepWhenConcurrencyLimitsExceeded)
			}
		}

		conn = nil
	}
}
