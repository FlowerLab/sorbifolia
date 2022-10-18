package strutils

import (
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

// Reverse an utf8 encoded string.
func Reverse(str string) string {
	b := pyrokinesis.String.ToBytes(str)
	bb := make([]byte, len(str))
	copy(bb, b)
	ReverseBytes(bb)
	return pyrokinesis.Bytes.ToString(bb)
}

// ReverseBytes an utf8 encoded string.
func ReverseBytes(str []byte) {
	n := len(str)
	if n < 1 {
		return
	}

	reverse(str)

	utf8Len := 0
	for i := 0; i < n; i++ {
		switch first[str[i]] {
		case xx:
			utf8Len++
		case as:
			utf8Len = 0
		case s1, s2, s3, s4, s5, s6, s7:
			if utf8Len >= 4 {
				utf8Len = 4
			} else {
				utf8Len++
			}
			char := str[i+1-utf8Len : i+1]
			charLen := len(char)

			reverse(char)
			if size := checkUTF8Len(char); size == 1 && charLen > 1 {
				reverse(char)
			}
			utf8Len = 0
		}
	}
}

func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
