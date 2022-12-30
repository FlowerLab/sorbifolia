package httpheader

import (
	"bytes"
)

type ContentDisposition []byte

func (v ContentDisposition) Type() []byte { return getHeaderValue(v) }

func (v ContentDisposition) Param(p []byte) (b []byte) {
	eachValueWithSemi(getHeaderParam(v), func(val []byte) bool {
		if k, value := parseKVWithEqual(val); bytes.EqualFold(k, p) {
			b = cleanQuotationMark(value)
			return false
		}
		return true
	})
	return
}

func (v ContentDisposition) Name() (b []byte)     { return v.Param([]byte("name")) }
func (v ContentDisposition) Filename() (b []byte) { return v.Param([]byte("filename")) }
