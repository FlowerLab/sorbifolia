//go:build goexperiment.arenas

package http

import (
	"net"
)

type DialFunc func(addr string) (net.Conn, error)

type Client struct {
	UserAgent []byte

	// Callback for establishing new connections to hosts.
	//
	// Default Dial is used if not set.
	Dial DialFunc
}

func (c Client) Do(req *Request, resp *Response) error {
	return nil
}

type ClientPool struct {
}
