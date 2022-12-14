package http

import (
	"net"
	"time"
)

type Context struct {
	c net.Conn

	Request  Request
	Response Response
}

func (c Context) Deadline() (deadline time.Time, ok bool) { panic("implement me") }
func (c Context) Done() <-chan struct{}                   { panic("implement me") }
func (c Context) Err() error                              { panic("implement me") }
func (c Context) Value(key any) any                       { panic("implement me") }
