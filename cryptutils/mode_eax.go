package cryptutils

import (
	cc "crypto/cipher"

	"go.x2ox.com/sorbifolia/cryptutils/cipher"
)

func EAX(block cc.Block) (AEAD, error) {
	return cipher.NewEAX(block)
}
