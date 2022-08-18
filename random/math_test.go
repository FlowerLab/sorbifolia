package random

import (
	"testing"
)

func TestMathRand(t *testing.T) {
	fr := NewMathRand()
	if len(fr.RandString(10)) != 10 {
		t.Error("1")
	}

	fr = fr.SetRandBytes([]byte("123456"))
	if len(fr.RandString(10)) != 10 {
		t.Error("test fail")
	}
}

func TestMathRandRepeatable(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	NewMathRand().SetRandBytes([]byte("11")).RandString(1)
}

func TestMathRandTooLong(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("test fail")
		}
	}()

	NewMathRand().SetRandBytes(make([]byte, 257)).RandString(1)
}
