package random

import (
	"testing"
)

func TestSafeRand(t *testing.T) {
	fr := NewSafeRand()
	if len(fr.RandString(10)) != 10 {
		t.Error("1")
	}

	fr = fr.SetRandBytes([]byte("123456"))
	if len(fr.RandString(10)) != 10 {
		t.Error("test fail")
	}
}

func TestSafeRandRepeatable(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	NewFastRand().SetRandBytes([]byte("11")).RandString(1)
}

func TestSafeRandTooLong(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	NewFastRand().SetRandBytes(make([]byte, 257)).RandString(1)
}
