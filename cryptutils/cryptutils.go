package cryptutils

import (
	"crypto/cipher"
)

type CryptStream interface {
	Encrypt(dst, src []byte)
	Decrypt(dst, src []byte)
}

type CryptBlock interface {
	Encrypt([]byte) []byte
	Decrypt([]byte) []byte
}

type AEAD interface {
	cipher.AEAD
}
