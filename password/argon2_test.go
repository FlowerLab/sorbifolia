package password

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"
)

func testCompare(t *testing.T, g Generator, password string) {
	hash := g.MustGenerate(password)
	if !g.Compare(hash, password) {
		t.Errorf("%s | %s not match", password, hash)
	}
	t.Logf("%s | %s", password, hash)
}

func TestArgon2(t *testing.T) {
	testCompare(t, New(), "password")
}

func TestArgon2_CompareFail(t *testing.T) {
	g := New()
	t.Run("", func(t *testing.T) {
		if g.Compare("", "") {
			t.Error("fail")
		}
	})

	t.Run("", func(t *testing.T) {
		if g.Compare(base64.RawStdEncoding.EncodeToString([]byte("1234567890")), "") {
			t.Error("fail")
		}
	})

	t.Run("", func(t *testing.T) {
		if g.Compare("AAAAAQABAAABAAAALZP/dD6HbO0SPK8Zijd/ivOT/3G3Wj1SKzrkIKs3REnw", "1") {
			t.Error("fail")
		}
	})
}

type errReader struct{}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("OEF")
}

func TestArgon2_MustGenerate(t *testing.T) {
	defer func() { _ = recover() }()

	rand.Reader = errReader{}
	g := New()
	g.MustGenerate("123456")

	t.Error("fail")
}
