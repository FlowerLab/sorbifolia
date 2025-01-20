package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
)

type Generator struct{}

func (g Generator) Ed25519() (ed25519.PrivateKey, error) {
	_, pri, err := ed25519.GenerateKey(nil)
	return pri, err
}

func (g Generator) ECDSA() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func (g Generator) RSA() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

func (g Generator) HMAC() ([]byte, error) {
	var arr [64]byte
	_, err := rand.Read(arr[:])
	return arr[:], err
}
