package strutils

import (
	"bytes"
	"testing"
)

func TestAppend(t *testing.T) {
	t.Parallel()

	if Append("asd", " def") != "asd def" {
		t.Error()
	}
}

func TestAppendToBytes(t *testing.T) {
	t.Parallel()

	if !bytes.Equal(AppendToBytes([]byte("asd"), " def"), []byte("asd def")) {
		t.Error()
	}

	b := make([]byte, 3, 10)
	copy(b, "asd")
	if !bytes.Equal(AppendToBytes(b, " def"), []byte("asd def")) {
		t.Error()
	}
}

func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Append("asd", " def")
	}
	b.ReportAllocs()
}

func BenchmarkAppendToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AppendToBytes([]byte("asd"), " def")
	}
	b.ReportAllocs()
}

func BenchmarkAppendToBytesNoCap(b *testing.B) {
	buf := make([]byte, 3, 3+b.N)
	copy(buf, "asd")

	for i := 0; i < b.N; i++ {
		AppendToBytes(buf, " ")
	}
	b.ReportAllocs()
}
