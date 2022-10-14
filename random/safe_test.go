package random

import (
	"testing"
)

func TestSafeRand(t *testing.T) {
	t.Parallel()

	fr := Safe()
	if len(fr.RandString(10)) != 10 {
		t.Error("1")
	}

	fr = fr.SetRandBytes([]byte("123456"))
	if len(fr.RandString(10)) != 10 {
		t.Error("test fail")
	}
}

func TestSafeRandRepeatable(t *testing.T) {
	t.Parallel()

	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	Safe().SetRandBytes([]byte("11")).RandString(1)
}

func TestSafeRandTooLong(t *testing.T) {
	t.Parallel()

	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	Safe().SetRandBytes(make([]byte, 257)).RandString(1)
}

func BenchmarkSafeRand(b *testing.B) {
	fr := Safe()
	for i := 0; i < b.N; i++ {
		fr.RandString(1123)
	}
}
