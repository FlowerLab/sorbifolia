package http

import (
	"context"
	"net"
	"sync"
	"time"

	"go.x2ox.com/sorbifolia/http/httpmessage"
)

type Context struct {
	c net.Conn
	s *Server

	id   uint64
	time time.Time
	addr net.Addr

	Request  httpmessage.Request
	Response httpmessage.Response
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
	c.Request.Reset()
	c.Response.Reset()
}

func (c *Context) Reset() {
	c.c = nil
	c.s = nil
	c.addr = nil
	c.Request.Reset()
	c.Response.Reset()
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
