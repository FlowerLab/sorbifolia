package httpconn

import (
	"crypto/tls"
)

type connTLSer interface {
	ConnectionState() tls.ConnectionState
}
