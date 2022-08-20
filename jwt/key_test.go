package jwt

import (
	"testing"
)

func TestGenerator(t *testing.T) {
	g := Generator{}

	if _, err := g.Ed25519(); err != nil {
		t.Error(err)
	}
	if _, err := g.ECDSA(); err != nil {
		t.Error(err)
	}
	if _, err := g.RSA(); err != nil {
		t.Error(err)
	}
	if _, err := g.HMAC(); err != nil {
		t.Error(err)
	}
}
