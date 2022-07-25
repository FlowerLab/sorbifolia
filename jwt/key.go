package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
)

type Generator struct {
}

func (g Generator) Ed25519() (ed25519.PrivateKey, error) {
	_, pri, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return pri, err
}

func (g Generator) ECDSA() (*ecdsa.PrivateKey, error) {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return pri, nil
}

func (g Generator) RSA() (*rsa.PrivateKey, error) {
	pri, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return pri, nil
}

func (g Generator) HMAC() ([]byte, error) {
	var arr [64]byte
	if _, err := rand.Read(arr[:]); err != nil {
		return nil, err
	}
	return arr[:], nil
}
