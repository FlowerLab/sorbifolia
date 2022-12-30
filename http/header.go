package http

import (
	"bytes"
	"crypto/tls"

	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/status"
)

type RequestHeader struct {
	KVs

	Accept           httpheader.Accept
	AcceptEncoding   httpheader.AcceptEncoding
	AcceptLanguage   httpheader.AcceptLanguage
	ContentLength    httpheader.ContentLength
	ContentType      httpheader.ContentType
	Cookie           httpheader.Cookie
	Host             httpheader.Host
	UserAgent        httpheader.UserAgent
	TransferEncoding httpheader.TransferEncoding

	// Trailer          httpheader.Trailer
	// TrailerHeader    KVs

	RemoteAddr []byte
	RequestURI []byte
	URL        URL
	TLS        *tls.ConnectionState
	Close      bool
}

func (rh *RequestHeader) RawParse() error {
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
		case bytes.EqualFold(kv.K, char.ContentLength):
			rh.ContentLength = kv.Val()
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
	KVs

	StatusCode      status.Status // e.g. 200
	ContentEncoding httpheader.ContentEncoding
	ContentLength   httpheader.ContentLength
	ContentType     httpheader.ContentType
	Date            httpheader.Date
	SetCookies      httpheader.SetCookies

	Close bool
}

func (rh *ResponseHeader) RawParse() error {
	rh.Each(func(kv KV) bool {
		switch {
		case bytes.EqualFold(kv.K, char.Connection):
			if bytes.EqualFold(kv.Val(), char.Close) {
				rh.Close = true
			}
		case bytes.EqualFold(kv.K, char.ContentLength):
			rh.ContentLength = kv.Val()
		case bytes.EqualFold(kv.K, char.SetCookie):
			rh.SetCookies = append(rh.SetCookies, kv.Val())
		}
		return true
	})

	return nil
}
