package httpconn

import (
	"crypto/tls"
	"net"
)

type TCP struct {
	ln     net.Listener
	tls    bool
	config *tls.Config
}

func newTCP(tln *net.TCPListener, tc *tls.Config) *TCP {
	return &TCP{ln: tln, tls: tc != nil, config: tc}
}

func (t *TCP) Addr() net.Addr { return t.ln.Addr() }
func (t *TCP) Accept() (conn net.Conn, err error) {
	if conn, err = t.ln.Accept(); err != nil || !t.tls {
		return
	}

	tc := tls.Server(conn, t.config)
	if err = tc.Handshake(); err != nil {
		return
	}

	return tc, nil
}

func (t *TCP) Close() error { return t.ln.Close() }
func (t *TCP) TLS() bool    { return t.tls }
