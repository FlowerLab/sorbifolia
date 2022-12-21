package http

import (
	"bytes"
	"crypto/tls"
	"time"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/status"
)

type RequestHeader struct {
	*KVs

	ContentType   []byte
	ContentLength int64
	Close         bool

	Accept         []byte
	AcceptEncoding []byte
	AcceptLanguage []byte
	UserAgent      []byte
	Cookie         []byte

	Host       []byte
	RemoteAddr []byte
	RequestURI []byte
	URL        *URL
	TLS        *tls.ConnectionState
}

func (rh *RequestHeader) init() error {
	rh.Each(func(kv KV) bool {
		switch {
		case bytes.EqualFold(kv.K, char.Accept):
			rh.Accept = kv.Val()
		case bytes.EqualFold(kv.K, char.AcceptEncoding):
			rh.AcceptEncoding = kv.Val()
		case bytes.EqualFold(kv.K, char.AcceptLanguage):
			rh.AcceptLanguage = kv.Val()
		case bytes.EqualFold(kv.K, char.Connection):
			if bytes.EqualFold(kv.Val(), char.Close) {
				rh.Close = true
			}
		case bytes.EqualFold(kv.K, char.Cookie):
			rh.Cookie = kv.Val()
		case bytes.EqualFold(kv.K, char.Host):
			rh.Host = kv.Val()
		case bytes.EqualFold(kv.K, char.UserAgent):
			rh.UserAgent = kv.Val()
		}
		return true
	})
	return rh.URL.Parse(rh.Host, rh.RequestURI, rh.TLS != nil)
}

type ResponseHeader struct {
	*KVs

	StatusCode status.Status // e.g. 200

	ContentType   []byte
	Date          time.Time
	ContentLength int64
	Close         bool
	Host          []byte
	RemoteAddr    []byte
	RequestURI    []byte
}
