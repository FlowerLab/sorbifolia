package httpheader

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
)

type ContentDisposition Value

func (v ContentDisposition) Type() []byte {
	if i := bytes.IndexByte(v, char.Semi); i >= 0 {
		return v[:i]
	}
	return v
}

func (v ContentDisposition) Param(p []byte) []byte {
	i := bytes.IndexByte(v, char.Semi)
	if i < 0 {
		return nil
	}
	b := v[i+1:]

	for {
		if len(b) == 0 {
			return nil
		}
		if b[0] == char.Space {
			b = b[1:]
			continue
		}

		if i = bytes.IndexByte(b, char.Semi); i == -1 {
			i = len(b) - 1
		}

		kv := b[:i]
		if len(b) < i+1 {
			i = len(b) - 2
		}
		b = b[i+1:]

		if i = bytes.IndexByte(kv, char.Equal); i == -1 {
			continue
		}
		if bytes.EqualFold(kv[:i], p) {
			return cleanQuotationMark(kv[i+1:])
		}
	}
}
