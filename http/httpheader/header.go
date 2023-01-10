package httpheader

import (
	"bytes"
	"strconv"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/kv"
)

type Header struct {
	kv.KVs
}

func (h *Header) Accept() Accept                 { return h.GetValue(char.Accept) }
func (h *Header) AcceptEncoding() AcceptEncoding { return h.GetValue(char.AcceptEncoding) }
func (h *Header) AcceptLanguage() AcceptLanguage { return h.GetValue(char.AcceptLanguage) }
func (h *Header) ContentLength() ContentLength   { return h.GetValue(char.ContentLength) }
func (h *Header) ContentType() ContentLength     { return h.GetValue(char.ContentType) }
func (h *Header) Cookie() ContentLength          { return h.GetValue(char.Cookie) }
func (h *Header) Host() Host                     { return h.GetValue(char.Host) }
func (h *Header) UserAgent() ContentLength       { return h.GetValue(char.UserAgent) }
func (h *Header) TransferEncoding() TransferEncoding {
	return h.GetValue(char.TransferEncoding)
}
func (h *Header) Trailer() Trailer { return h.GetValue(char.Trailer) }

func (h *Header) SetContentLength(i int)  { h.setI(char.ContentLength, i) }
func (h *Header) SetContentType(b []byte) { h.Set(char.ContentType, b) }
func (h *Header) SetHost(b []byte)        { h.Set(char.Host, b) }

// func (h *Header) setS(k []byte, v string) {
// 	val := h.GetOrAdd(k)
// 	val.V = append(val.V, v...)
// }

func (h *Header) setI(k []byte, i int) {
	v := h.GetOrAdd(k)
	v.V = strconv.AppendInt(v.V, int64(i), 10)
}

func (h *Header) Parse(b []byte) {
	if len(b) == 0 {
		return
	}
	h.KVs.PreAlloc(bytes.Count(b, char.CRLF) + 1)

	for len(b) > 0 {
		idx := bytes.Index(b, char.CRLF)
		switch idx {
		case 0:
		case -1:
			h.KVs.AddHeader(b)
			idx = len(b) - 2
		default:
			h.KVs.AddHeader(b[:idx])
		}

		b = b[idx+2:]
	}
}

func (h *Header) Reset() { h.KVs.Reset() }

func ParseHeader(b []byte) (k, v []byte, null bool) {
	idx := bytes.IndexByte(b, char.Colon)
	if idx == -1 {
		k, null = b, true
		return
	}

	k = b[:idx]
	idx++
	for ; idx < len(b); idx++ {
		if b[idx] != char.Space {
			v = b[idx:]
			break
		}
	}

	return
}

func AppendHeader(dst []byte, v kv.KV) []byte {
	if dst = append(dst, v.K...); !v.Null {
		dst = append(dst, char.Colon)
		dst = append(dst, char.Space)
		dst = append(dst, v.V...)
	}
	dst = append(dst, char.CRLF...)
	return dst
}

func AppendHeaders(dst []byte, v kv.KVs) []byte {
	v.Each(func(kv kv.KV) bool { dst = AppendHeader(dst, kv); return true })
	dst = append(dst, char.CRLF...)
	return dst
}
