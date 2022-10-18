package strutils

import (
	"testing"
)

var checkUTF8LenTests = []struct {
	b    []byte
	size int
}{
	{[]byte{'a'}, 1},                  // as
	{[]byte("é"), 2},                  // s1
	{[]byte("我"), 3},                  // s3
	{[]byte(string('\U0010FFFF')), 4}, // s7
	{[]byte{0xf4, 0b10000100, 0b10000101, 0b10000110}, 4}, // s7

	{[]byte{0xe0, 0b10110001, 0b10100001}, 3},             // s2
	{[]byte{0xed, 0b10011111, 0b10000111}, 3},             // s4
	{[]byte{0xf0, 0b10010000, 0b10011000, 0b10011100}, 4}, // s5
	{[]byte{0xf1, 0b10000010, 0b10100000, 0b10101010}, 4}, // s6

	{[]byte{0xa9}, 1}, // xx

	{[]byte{}, 0}, // line 6

	// add test for p[1] < accept.lo || accept.hi < p[1]
	{[]byte{0xc2, 0b00110001}, 1}, // line 22
	{[]byte{0xc2, 0b11111111}, 1}, // line 22

	{[]byte{0xe0, 0b10110001, 0b00110001}, 1}, // line 29
	{[]byte{0xe0, 0b10110001, 0b11111111}, 1}, // line 29

	{[]byte{0xf1, 0b10110001, 0b10011000, 0b00110001}, 1}, // line 36
	{[]byte{0xf1, 0b10110001, 0b10011000, 0b11111111}, 1}, // line 36

	{[]byte{0xf1, 0b10110001, 0b10011000}, 1},
}

func Test_checkUTF8Len(t *testing.T) {
	t.Parallel()

	for _, v := range checkUTF8LenTests {
		if checkUTF8Len(v.b) != v.size {
			t.Error("fail")
		}
	}
}
