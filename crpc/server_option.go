package crpc

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/net/http2"
)

type ServerOption struct {
	h2c  *http2.Server
	srv  *http.Server
	cert *tls.Certificate

	addr   string
	handle HttpHandle

	ham *healthAndMetrics
}

type ApplyToServer func(*ServerOption)

func WithH2C(h2c *http2.Server) ApplyToServer {
	if h2c == nil {
		h2c = &http2.Server{}
	}
	return func(o *ServerOption) { o.h2c = h2c }
}

func WithCertFile(certFile, keyFile string) ApplyToServer {
	val, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	return func(o *ServerOption) { o.cert = &val }
}

func WithCertPEM(cert, key []byte) ApplyToServer {
	val, err := tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	return func(o *ServerOption) { o.cert = &val }
}

func WithHTTPServer(s *http.Server) ApplyToServer {
	if s == nil {
		s = &http.Server{}
	}
	return func(o *ServerOption) { o.srv = s }
}

func WithHandle(h HttpHandle) ApplyToServer {
	return func(o *ServerOption) { o.handle = h }
}

func WithAddr(addr string) ApplyToServer {
	return func(o *ServerOption) { o.addr = addr }
}

func WithHealthAndMetrics(addr, _path string) ApplyToServer {
	return func(o *ServerOption) { o.ham = &healthAndMetrics{addr: addr, path: _path} }
}
