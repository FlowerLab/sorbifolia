package httpheader

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/kv"
)

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
