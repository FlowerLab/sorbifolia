package cryptutils

import (
	"crypto/cipher"
)

func GCM(block cipher.Block) (AEAD, error) {
	return cipher.NewGCM(block)
}
