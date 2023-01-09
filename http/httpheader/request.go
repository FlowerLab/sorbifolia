package httpheader

import (
	"bytes"
	"crypto/tls"
	"strconv"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/kv"
	"go.x2ox.com/sorbifolia/http/url"
)

type RequestHeader struct {
	kv.KVs

	RemoteAddr []byte
	RequestURI []byte
	URL        url.URL
	TLS        *tls.ConnectionState
	Close      bool
}

func (rh *RequestHeader) Accept() Accept                 { return rh.GetValue(char.Accept) }
func (rh *RequestHeader) AcceptEncoding() AcceptEncoding { return rh.GetValue(char.AcceptEncoding) }
func (rh *RequestHeader) AcceptLanguage() AcceptLanguage { return rh.GetValue(char.AcceptLanguage) }
func (rh *RequestHeader) ContentLength() ContentLength   { return rh.GetValue(char.ContentLength) }
func (rh *RequestHeader) ContentType() ContentLength     { return rh.GetValue(char.ContentType) }
func (rh *RequestHeader) Cookie() ContentLength          { return rh.GetValue(char.Cookie) }
func (rh *RequestHeader) Host() Host                     { return rh.GetValue(char.Host) }
func (rh *RequestHeader) UserAgent() ContentLength       { return rh.GetValue(char.UserAgent) }
func (rh *RequestHeader) TransferEncoding() TransferEncoding {
	return rh.GetValue(char.TransferEncoding)
}
func (rh *RequestHeader) Trailer() Trailer { return rh.GetValue(char.Trailer) }

func (rh *RequestHeader) SetContentLength(i int)  { rh.setI(char.ContentLength, i) }
func (rh *RequestHeader) SetContentType(b []byte) { rh.Set(char.ContentType, b) }
func (rh *RequestHeader) SetHost(b []byte)        { rh.Set(char.Host, b) }

func (rh *RequestHeader) setS(k []byte, v string) {
	val := rh.GetOrAdd(k)
	val.V = append(val.V, v...)
}

func (rh *RequestHeader) setI(k []byte, i int) {
	v := rh.GetOrAdd(k)
	v.V = strconv.AppendInt(v.V, int64(i), 10)
}

func (rh *RequestHeader) Parse(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	rh.KVs.PreAlloc(bytes.Count(b, char.CRLF) + 1)

	for len(b) > 0 {
		idx := bytes.Index(b, char.CRLF)
		switch idx {
		case 0:
		case -1:
			rh.KVs.AddHeader(b)
			idx = len(b) - 2
		default:
			rh.KVs.AddHeader(b[:idx])
		}

		b = b[idx+2:]
	}

	return rh.URL.Parse(rh.KVs.Get(char.Host).V, rh.RequestURI, rh.TLS != nil)
}

func (rh *RequestHeader) Reset() {
	rh.KVs.Reset()
	rh.URL.Reset()
	rh.TLS = nil
	rh.Close = false
}
