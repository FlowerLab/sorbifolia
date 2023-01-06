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

func (rh *RequestHeader) Reset() {
	rh.KVs.Reset()
	rh.Accept = rh.Accept[:0]
	rh.AcceptEncoding = rh.AcceptEncoding[:0]
	rh.AcceptLanguage = rh.AcceptLanguage[:0]
	rh.ContentLength = rh.ContentLength[:0]
	rh.ContentType = rh.ContentType[:0]
	rh.Cookie = rh.Cookie[:0]
	rh.Host = rh.Host[:0]
	rh.UserAgent = rh.UserAgent[:0]
	rh.TransferEncoding = rh.TransferEncoding[:0]
	rh.Trailer = rh.Trailer[:0]
	rh.TrailerHeader.Reset()
	rh.RemoteAddr = rh.RemoteAddr[:0]
	rh.RequestURI = rh.RequestURI[:0]
	rh.URL.Reset()
	rh.TLS = nil
	rh.Close = false
}
