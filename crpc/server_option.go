package crpc

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
)

type ServerOption struct {
	h2c  *http2.Server
	srv  *http.Server
	cert *tls.Certificate
	cors func(h http.Handler) http.Handler

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

func WithCert(cert, key string) ApplyToServer {
	return WithCertPEM([]byte(cert), []byte(key))
}

func WithCertFromCheck(env string, path ...string) ApplyToServer {
	if env != "" && os.Getenv(fmt.Sprintf("%s_CRT", env)) != "" {
		return WithCert(
			os.Getenv(fmt.Sprintf("%s_CRT", env)),
			os.Getenv(fmt.Sprintf("%s_KEY", env)),
		)
	}
	for _, v := range path {
		if _, err := os.Stat(v + ".crt"); err == nil {
			return WithCertFile(v+".crt", v+".key")
		}
	}

	panic("env and path not found")
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

func WithCORS(fn func(h http.Handler) http.Handler) ApplyToServer {
	if fn == nil {
		fn = cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{"Authorization", "Content-Type", "Connect-Protocol-Version", "Connect-Timeout-Ms", "Grpc-Timeout", "X-Grpc-Web", "X-User-Agent"},
			ExposedHeaders: []string{"Grpc-Status", "Grpc-Message", "Grpc-Status-Details-Bin"},
		}).Handler
	}
	return func(o *ServerOption) { o.cors = fn }
}
