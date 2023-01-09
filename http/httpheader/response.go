package httpheader

import (
	"strconv"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/kv"
)

type ResponseHeader struct {
	kv.KVs

	Close bool
}

func (rh *ResponseHeader) Reset() {
	rh.KVs.Reset()
}

func (rh *ResponseHeader) ContentLength() ContentLength { return rh.GetValue(char.ContentLength) }
func (rh *ResponseHeader) ContentType() ContentLength   { return rh.GetValue(char.ContentType) }
func (rh *ResponseHeader) Cookie() ContentLength        { return rh.GetValue(char.Cookie) }
func (rh *ResponseHeader) Trailer() Trailer             { return rh.GetValue(char.Trailer) }

func (rh *ResponseHeader) SetContentLength(i int)  { rh.setI(char.ContentLength, i) }
func (rh *ResponseHeader) SetContentType(b []byte) { rh.Set(char.ContentLength, b) }

func (rh *ResponseHeader) setS(k []byte, v string) {
	val := rh.GetOrAdd(k)
	val.V = append(val.V, v...)
}

func (rh *ResponseHeader) setI(k []byte, i int) {
	v := rh.GetOrAdd(k)
	v.V = strconv.AppendInt(v.V, int64(i), 10)
}
