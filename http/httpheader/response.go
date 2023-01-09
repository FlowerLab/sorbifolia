package httpheader

import (
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
	return nil
}

func (rh *ResponseHeader) Reset() {
	rh.KVs.Reset()
	rh.ContentLength = rh.ContentLength[:0]
	rh.ContentType = rh.ContentType[:0]
	rh.Close = false

	for i := range rh.SetCookies {
		rh.SetCookies[i] = rh.SetCookies[i][:0]
	}
	rh.SetCookies = rh.SetCookies[:0]
}
