package httpheader

import (
	"crypto/tls"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/url"
)

type RequestHeader struct {
	Header

	RemoteAddr []byte
	RequestURI []byte
	URL        url.URL
	TLS        *tls.ConnectionState
	Close      bool
}

func (rh *RequestHeader) Parse(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	rh.Header.Parse(b)
	return rh.URL.Parse(rh.Header.Get(char.Host).V, rh.RequestURI, rh.TLS != nil)
}

func (rh *RequestHeader) Reset() {
	rh.Header.Reset()
	rh.URL.Reset()
	rh.TLS = nil
	rh.Close = false
}
