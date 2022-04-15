package password

import (
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
