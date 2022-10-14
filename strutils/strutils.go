package coarsetime

import (
	"unicode/utf8"

	"go.x2ox.com/sorbifolia/pyrokinesis"
)

// ReverseA an utf8 encoded string.
func ReverseA(str string) string {
	b := pyrokinesis.String.ToBytes(str)
	// s := pyrokinesis.Bytes.ToString(str)

	// tail := len(str)

	for len(str) > 0 {
		utf8.DecodeRune(b)
		// _, size = utf8.DecodeRuneInString(s)
		// tail -= size
		// buf = append(buf[:tail], []byte(str[:size])...)
		// str = str[size:]
	}

	return "buf"
}

// Reverse an utf8 encoded string.
func Reverse(str string) string {
	r := []rune(str)

	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// ReverseB a utf8 encoded string.
func ReverseB(str string) string {
	var size int

	tail := len(str)
	buf := make([]byte, tail)
	s := buf

	for len(str) > 0 {
		_, size = utf8.DecodeRuneInString(str)
		tail -= size
		s = append(s[:tail], []byte(str[:size])...)
		str = str[size:]
	}

	return string(buf)
}

// ReverseA1 an utf8 encoded string.
func ReverseA1(str []byte) []byte {
	var size int
	s := pyrokinesis.Bytes.ToString(str)

	tail := len(str)
	buf := make([]byte, tail)

	for len(str) > 0 {
		_, size = utf8.DecodeRuneInString(s)
		tail -= size
		buf = append(buf[:tail], []byte(str[:size])...)
		str = str[size:]
	}

	return buf
}
