//go:build goexperiment.arenas

package http

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
)

func TestS(t *testing.T) {
	s := &Server{
		Name:                  []byte("aa"),
		MaxRequestHeaderSize:  defaultMaxRequestHeaderSize,
		MaxRequestBodySize:    defaultMaxRequestBodySize,
		StreamRequestBodySize: defaultMaxRequestBodySize,

		Handler: func(ctx *Context) {},
	}

	ln, _ := net.Listen("tcp", "127.0.0.1:8808")
	s.Serve(ln)
}

func TestK(t *testing.T) {
	as(1)
	as("11")
	as(bytes.NewReader([]byte("a")))
}
func as(a any) {
	switch a := a.(type) {
	case string:
		fmt.Println("s", a)
	case int64, int:
		fmt.Println("i")
	case io.Reader:
		fmt.Println("ir")
	}
}
