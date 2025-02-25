package password

import (
	"encoding/base64"
	"testing"
)

func TestMustGenerate(t *testing.T) {
	t.Parallel()

	hash := MustGenerate("password")
	if !Compare(hash, "password") {
		t.Errorf("%s | %s not match", "password", hash)
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	hash, _ := Generate("password")
	if !Compare(hash, "password") {
		t.Errorf("%s | %s not match", "password", hash)
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	hash := MustGenerate("password")
	if !Compare(hash, "password") {
		t.Errorf("%s | %s not match", "password", hash)
	}
}

func TestFail(t *testing.T) {
	t.Parallel()

	t.Run("", func(t *testing.T) {
		if Compare("", "") {
			t.Error("fail")
		}
	})

	t.Run("", func(t *testing.T) {
		if Compare(base64.RawStdEncoding.EncodeToString([]byte("1234567890")), "") {
			t.Error("fail")
		}
	})

	t.Run("", func(t *testing.T) {
		if Compare("AAAAAQABAAABAAAALZP/dD6HbO0SPK8Zijd/ivOT/3G3Wj1SKzrkIKs3REnw", "1") {
			t.Error("fail")
		}
	})
}
