package strutils

import (
	"bytes"
	"testing"
)

var (
	testStr1a = "阿ab12三啊11232123实打123实121233312打3算123"
	testStr1b = "321算3打213332121实321打实32123211啊三21ba阿"
)

func TestReverse(t *testing.T) {
	t.Parallel()

	s := Reverse(testStr1a)
	if s != testStr1b {
		t.Error("fail")
	}
}

func TestReverseBytes(t *testing.T) {
	t.Parallel()

	b := []byte(testStr1a)
	ReverseBytes(b)
	if !bytes.Equal(b, []byte(testStr1b)) {
		t.Error("fail")
	}
}

func Test_ReverseBytesEmpty(t *testing.T) {
	t.Parallel()
	var b []byte
	ReverseBytes(b)
	if !bytes.Equal(b, []byte{}) {
		t.Error("fail")
	}
}

func TestBoundary(t *testing.T) {
	t.Parallel()
	b := []byte{0xe0, 0b10110001, 0b00110001}
	ReverseBytes(b)
	if !bytes.Equal(b, []byte{0b00110001, 0b10110001, 0xe0}) {
		t.Error("fail")
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
