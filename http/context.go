package http

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"go.x2ox.com/sorbifolia/http/httpbody"
)

type Context struct {
	c net.Conn
	s *Server

	id   uint64
	time time.Time
	addr net.Addr

	Request  Request
	Response Response
}

func (c *Context) Deadline() (deadline time.Time, ok bool) { return }
func (c *Context) Done() <-chan struct{}                   { return c.s.done }
func (c *Context) Err() error {
	select {
	case <-c.s.done:
		return context.Canceled
	default:
		return nil
	}
}
func (c *Context) Value(key any) any {
	panic("implement me")
}

func (c *Context) cleanup() {
	c.Request.Header.Reset()
	c.Response.Header.Reset()

	if c.Request.Body != nil {
		_ = c.Request.Body.Close()
		if p, ok := c.Request.Body.(httpbody.Pool); ok {
			httpbody.Release(p)
		}
		c.Request.Body = nil
	}

	if c.Response.Body != nil {
		if bc, ok := c.Response.Body.(io.Closer); ok {
			_ = bc.Close()
		}
		if p, ok := c.Response.Body.(httpbody.Pool); ok {
			httpbody.Release(p)
		}
		c.Response.Body = nil
	}
}

func (c *Context) Reset() {
	c.c = nil
	c.s = nil
	c.addr = nil
	c.Request.Header.Reset()
	c.Response.Header.Reset()

	if c.Request.Body != nil {
		_ = c.Request.Body.Close()
		if p, ok := c.Request.Body.(httpbody.Pool); ok {
			httpbody.Release(p)
		}
		c.Request.Body = nil
	}

	if c.Response.Body != nil {
		if bc, ok := c.Response.Body.(io.Closer); ok {
			_ = bc.Close()
		}
		if p, ok := c.Response.Body.(httpbody.Pool); ok {
			httpbody.Release(p)
		}
		c.Response.Body = nil
	}
}

func AcquireContext() *Context {
	if v := _ContextPool.Get(); v != nil {
		return v.(*Context)
	}
	return &Context{}
}

func ReleaseContext(c *Context) {
	c.Reset()
	_ContextPool.Put(c)
}

var (
	_ContextPool = sync.Pool{}
)
