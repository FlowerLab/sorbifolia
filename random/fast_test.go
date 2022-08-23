package random

import (
	"testing"
)

func TestNewFastRand(t *testing.T) {
	fr := NewFastRand()
	if len(fr.RandString(10)) != 10 {
		t.Error("1")
	}

	fr = fr.SetRandBytes([]byte("123456"))
	if len(fr.RandString(10)) != 10 {
		t.Error("test fail")
	}
}

func TestNewFastRandRepeatable(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	NewFastRand().SetRandBytes([]byte("11")).RandString(1)
}

func TestNewFastRandTooLong(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	NewFastRand().SetRandBytes(make([]byte, 257)).RandString(1)
}

func BenchmarkFastRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFastRand().RandString(1123)
	}
}

func BenchmarkFastRand2(b *testing.B) {
	fr := NewFastRand()
	for i := 0; i < b.N; i++ {
		fr.RandString(1123)
	}
}

func TestNewFastRandTooLon(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	NewFastRand().SetRandBytes(make([]byte, 257)).RandString(1)
}

func TestFastRand64(t *testing.T) {
	fastRand64()
}
