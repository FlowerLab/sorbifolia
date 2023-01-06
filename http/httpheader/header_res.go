package httpheader

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/kv"
)

type ResponseHeader struct {
	kv.KVs

	ContentLength ContentLength
	ContentType   ContentType
	SetCookies    SetCookies

	Close bool
}

func (rh *ResponseHeader) RawParse() error {
	rh.Each(func(kv kv.KV) bool {
		switch {
		case bytes.EqualFold(kv.K, char.Connection):
			if bytes.EqualFold(kv.V, char.Close) {
				rh.Close = true
			}
		case bytes.EqualFold(kv.K, char.ContentLength):
			rh.ContentLength = kv.V
		case bytes.EqualFold(kv.K, char.SetCookie):
			rh.SetCookies = append(rh.SetCookies, kv.V)
		}
		return true
	})

	return nil
}
