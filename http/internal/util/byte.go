//go:generate go run byte_table_gen.go

package util

import (
	"bytes"
)

const (
	upperhex = "0123456789ABCDEF"
	lowerhex = "0123456789abcdef"
)

func ToLower(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		p := &b[i]
		*p = toLowerTable[*p]
	}
	return b
}

func ToUpper(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		p := &b[i]
		*p = toUpperTable[*p]
	}
	return b
}

func AppendQuotedPath(dst, src []byte) []byte {
	for _, c := range src {
		if quotedPathShouldEscapeTable[int(c)] != 0 {
			dst = append(dst, '%', upperhex[c>>4], upperhex[c&0xf])
		} else {
			dst = append(dst, c)
		}
	}
	return dst
}

// DecodeArgAppendNoPlus is almost identical to decodeArgAppend, but it doesn't
// substitute '+' with ' '.
//
// The function is copy-pasted from decodeArgAppend due to the performance
// reasons only.
func DecodeArgAppendNoPlus(dst, src []byte) []byte {
	idx := bytes.IndexByte(src, '%')
	if idx < 0 {
		// fast path: src doesn't contain encoded chars
		return append(dst, src...)
	} else {
		dst = append(dst, src[:idx]...)
	}

	// slow path
	for i := idx; i < len(src); i++ {
		c := src[i]
		if c == '%' {
			if i+2 >= len(src) {
				return append(dst, src[i:]...)
			}
			x2 := hex2intTable[src[i+2]]
			x1 := hex2intTable[src[i+1]]
			if x1 == 16 || x2 == 16 {
				dst = append(dst, '%')
			} else {
				dst = append(dst, x1<<4|x2)
				i += 2
			}
		} else {
			dst = append(dst, c)
		}
	}
	return dst
}
