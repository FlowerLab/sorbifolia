//go:build goexperiment.arenas

package http

import (
	"arena"
	"context"
	"net"
	"time"
)

type Context struct {
	a *arena.Arena
	c net.Conn
	s *Server

	id   uint64
	time time.Time
	addr net.Addr

	Request  Request
	Response Response

	robbery bool
}

// Robbery will make the Server never Free it, Must call Free() yourself later.
func (c *Context) Robbery() {
	if c.robbery {
		panic("Robbery() cannot be called multiple times")
	}
	c.robbery = true
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
