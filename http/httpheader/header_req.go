package httpheader

import (
	"bytes"
	"crypto/tls"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/kv"
	"go.x2ox.com/sorbifolia/http/url"
)

type RequestHeader struct {
	kv.KVs

	Accept           Accept
	AcceptEncoding   AcceptEncoding
	AcceptLanguage   AcceptLanguage
	ContentLength    ContentLength
	ContentType      ContentType
	Cookie           Cookie
	Host             Host
	UserAgent        UserAgent
	TransferEncoding TransferEncoding

	Trailer       Trailer
	TrailerHeader kv.KVs

	RemoteAddr []byte
	RequestURI []byte
	URL        url.URL
	TLS        *tls.ConnectionState
	Close      bool
}

func (rh *RequestHeader) RawParse() error {
	rh.Each(func(kv kv.KV) bool {
		switch {
		case bytes.EqualFold(kv.K, char.Accept):
			rh.Accept = kv.V
		case bytes.EqualFold(kv.K, char.AcceptEncoding):
			rh.AcceptEncoding = kv.V
		case bytes.EqualFold(kv.K, char.AcceptLanguage):
			rh.AcceptLanguage = kv.V
		case bytes.EqualFold(kv.K, char.Connection):
			if bytes.EqualFold(kv.V, char.Close) {
				rh.Close = true
			}
		case bytes.EqualFold(kv.K, char.ContentLength):
			rh.ContentLength = kv.V
		case bytes.EqualFold(kv.K, char.Cookie):
			rh.Cookie = kv.V
		case bytes.EqualFold(kv.K, char.Host):
			rh.Host = kv.V
		case bytes.EqualFold(kv.K, char.UserAgent):
			rh.UserAgent = kv.V
		}
		return true
	})

	return rh.URL.Parse(rh.Host, rh.RequestURI, rh.TLS != nil)
}
