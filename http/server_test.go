//go:build goexperiment.arenas

package http

import (
	"net"
	"testing"
)

func TestS(t *testing.T) {
	s := &Server{
		MaxRequestHeaderSize:  defaultMaxRequestHeaderSize,
		MaxRequestBodySize:    defaultMaxRequestBodySize,
		StreamRequestBodySize: defaultMaxRequestBodySize,

		Handler: func(ctx *Context) {},
	}

	ln, _ := net.Listen("tcp", "127.0.0.1:8808")
	s.Serve(ln)
}
