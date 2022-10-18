package strutils

import (
	"bytes"
	"testing"
)

var (
	testStr1a = "阿ab12三啊11232123实打123实121233312打3算123"
	testStr1b = "321算3打213332121实321打实32123211啊三21ba阿"
)

var reverseTestsStr = []struct {
	b       []byte
	reverse []byte
}{
	{[]byte(testStr1a), []byte(testStr1b)},
	{[]byte{}, []byte{}}, // empty
	{[]byte{0xe0, 0b10110001, 0b00110001}, []byte{0b00110001, 0b10110001, 0xe0}}, // boundary

	{
		[]byte{0xf1, 0b10110001, 0b10110001, 0b10110001, 0b10110001, 0b10110001, 0b00110001},
		[]byte{0b00110001, 0b10110001, 0b10110001, 0b10110001, 0b10110001, 0b10110001, 0xf1},
	},
	{
		[]byte{0xf1, 0b00110011, 0b00110001, 0b10110001},
		[]byte{0b10110001, 0b00110001, 0b00110011, 0xf1},
	},
	{
		[]byte{0xf1, 0b10110001, 0b10110001, 0b10110001, 0b00110001, 0b00110001, 0b10110001},
		[]byte{0b10110001, 0b00110001, 0b00110001, 0b10110001, 0b10110001, 0b10110001, 0xf1},
	},
}

func TestReverse(t *testing.T) {
	t.Parallel()

	s := Reverse(testStr1a)
	if s != testStr1b {
		t.Error("fail")
	}
}

func TestReverseBytes(t *testing.T) {
	t.Parallel()

	for _, v := range reverseTestsStr {
		ReverseBytes(v.b)
		if !bytes.Equal(v.b, v.reverse) {
			t.Error("fail")
		}
	}
}

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reverse(testStr1a)
	}
	b.ReportAllocs()
}

func BenchmarkReverseBytes(b *testing.B) {
	p := []byte(testStr1a)
	for i := 0; i < b.N; i++ {
		ReverseBytes(p)
	}
	b.ReportAllocs()
}
