package httpconn

import (
	"crypto/tls"
	"net"
)

type Listener interface {
	net.Listener
	TLS() bool
}

func NewTCP(tln *net.TCPListener, tc *tls.Config) *TCP { return &TCP{ln: tln, config: tc} }
func NewUDP(conn *net.UDPConn, tc *tls.Config) *UDP    { return &UDP{conn: conn, config: tc} }

type TCP struct {
	ln     *net.TCPListener
	config *tls.Config
}

type UDP struct {
	conn   *net.UDPConn
	config *tls.Config
}

func (t *TCP) Accept() (conn net.Conn, err error) {
	if conn, err = t.ln.Accept(); err != nil || !t.TLS() {
		return
	}

	tc := tls.Server(conn, t.config)
	if err = tc.Handshake(); err != nil {
		return
	}

	return tc, nil
}
func (t *TCP) Addr() net.Addr { return t.ln.Addr() }
func (t *TCP) Close() error   { return t.ln.Close() }
func (t *TCP) TLS() bool      { return t.config != nil }
func (t *UDP) Addr() net.Addr { return t.conn.LocalAddr() }
func (t *UDP) Close() error   { return t.conn.Close() }
func (t *UDP) TLS() bool      { return t.config != nil }
func (t *UDP) Accept() (conn net.Conn, err error) {
	panic("")
}

var (
	_ Listener = (*UDP)(nil)
	_ Listener = (*TCP)(nil)
)
