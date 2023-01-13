package httpconn

import (
	"net"

	"go.x2ox.com/sorbifolia/http/version"
)

type HTTPConn interface {
	net.Conn
	Version() version.Version
	TLS() bool
}

type Conn struct {
	net.Conn
	ver version.Version
}

func (c *Conn) Version() version.Version { return c.ver }
func (c *Conn) TLS() (ok bool)           { _, ok = c.Conn.(connTLSer); return }

var (
	_h2          = version.Version{Major: 2}
	_   HTTPConn = (*Conn)(nil)
)
