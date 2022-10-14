package password

import (
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
