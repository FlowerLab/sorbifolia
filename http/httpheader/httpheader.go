package httpheader

import (
	"go.x2ox.com/sorbifolia/http/internal/char"
)

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
