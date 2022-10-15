package strutils

import (
	"testing"
	"unicode/utf8"
)

var checkUTF8LenTests = []struct {
	b    []byte
	size int
}{
	{[]byte{'a'}, 1},
	{[]byte("é"), 2},
	{[]byte("我"), 3},
	{[]byte(string('\U0010FFFF')), 4},
}

func Test_checkUTF8Len(t *testing.T) {
	t.Parallel()

	for _, v := range checkUTF8LenTests {
		_, size := utf8.DecodeRune(v.b)
		if checkUTF8Len(v.b) != size {
			t.Error("fail")
		}
	}
}
