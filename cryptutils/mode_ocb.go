package cryptutils

import (
	cc "crypto/cipher"

	"go.x2ox.com/sorbifolia/cryptutils/cipher"
)

func OCB(block cc.Block) (AEAD, error) {
	return cipher.NewOCB(block)
}
