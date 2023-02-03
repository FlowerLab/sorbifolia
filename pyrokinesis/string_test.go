package pyrokinesis

import (
	"bytes"
	"testing"
)

func TestString_Copy(t *testing.T) {
	t.Parallel()

	s := "hello"
	if String.Copy(s) != s {
		t.Error("fail")
	}
}

func TestString_ToBytes(t *testing.T) {
	t.Parallel()

	s := "hello"
	if !bytes.Equal(String.ToBytes(s), []byte("hello")) {
		t.Error("fail")
	}

	if String.ToBytes("") != nil {
		t.Error("fail")
	}
}

func BenchmarkString_Copy(b *testing.B) {
	s := "hello"
	for i := 0; i < b.N; i++ {
		String.Copy(s)
	}
}

func BenchmarkString_ToBytes(b *testing.B) {
	s := "hello"
	for i := 0; i < b.N; i++ {
		String.ToBytes(s)
	}
}
