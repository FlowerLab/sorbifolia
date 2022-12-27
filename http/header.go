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

	ContentType   ContentType
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

	ContentType   ContentType
	Date          time.Time
	ContentLength int64
	Close         bool
	Host          []byte
	RemoteAddr    []byte
	RequestURI    []byte
}

type ContentType []byte

func (ct ContentType) MIME() []byte {
	if i := bytes.IndexByte(ct, char.Semi); i >= 0 {
		return ct[:i]
	}
	return ct
}

func (ct ContentType) Charset() []byte {
	var charset []byte = nil

	if i := bytes.Index(ct, char.Charset); i >= 0 {
		charset = ct[i+len(char.Charset):]
	}
	if len(charset) == 0 || charset[0] != char.Equal {
		return nil
	}
	charset = charset[1:]

	if i := bytes.IndexByte(charset, char.Semi); i >= 0 {
		charset = charset[:i]
	}

	return cleanQuotationMark(cleanTrailingSpaces(charset))
}

func (ct ContentType) Boundary() []byte {
	var boundary []byte = nil

	if i := bytes.Index(ct, char.Boundary); i >= 0 {
		boundary = ct[i+len(char.Boundary):]
	}
	if len(boundary) == 0 || boundary[0] != char.Equal {
		return nil
	}
	boundary = boundary[1:]

	if i := bytes.IndexByte(boundary, char.Semi); i >= 0 {
		boundary = boundary[:i]
	}

	return cleanTrailingSpaces(boundary)
}

func cleanTrailingSpaces(b []byte) []byte {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != char.Space {
			return b[:i]
		}
	}
	return nil
}

func cleanQuotationMark(b []byte) []byte {
	if len(b) < 2 || b[0] != char.QuotationMark || b[len(b)-1] != char.QuotationMark {
		return b
	}
	return b[1 : len(b)-1]
}

// QuotationMark
