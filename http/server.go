package http

import (
	"net"
)

type Handler func(ctx *Context)

type Server struct {
}

func (s *Server) Listen()              {}
func (s *Server) handle(conn net.Conn) {}

func (s *Server) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		s.handle(conn)
	}
}
