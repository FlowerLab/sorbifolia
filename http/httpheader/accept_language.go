package httpheader

import (
	"bytes"
	"strconv"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type AcceptLanguage Value

// Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5
// Set-Cookie: UserID=JohnDoe; Max-Age=3600; Version=1
// Content-Type: text/html; charset=utf-8
// Content-Disposition: attachment; filename="name.ext"

func (v AcceptLanguage) Each(fn func(value QualityValue) bool) {
	b := v

	for {
		if len(b) == 0 {
			return
		}

		i := bytes.IndexByte(b, char.Comma)
		if i < 0 {
			i = len(b)
		}

		val := b[:i]

		{
			if len(b) < i+1 {
				i = len(b) - 1
			}
			b = b[i+1:]

			j := 0
			for ; j < len(val); j++ {
				if val[j] != char.Space {
					break
				}
			}
			val = val[j:]
		}

		var qv QualityValue
		if i = bytes.IndexByte(val, char.Semi); i < 0 {
			qv.Value = val
			qv.Priority = 1
		} else {
			qv.Value = val[:i]
			val = val[i:]
			if i = bytes.IndexByte(val, char.Equal); i > 0 {
				qv.Priority, _ = strconv.ParseFloat(pyrokinesis.Bytes.ToString(val[i+1:]), 64)
			}
		}

		if !fn(qv) {
			return
		}
	}
}

type QualityValue struct {
	Value    []byte
	Priority float64 // 1.00 - 0.00
}
