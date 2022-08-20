package password

import (
	"encoding/base64"
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
