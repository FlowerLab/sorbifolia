package cryptutils

import (
	cc "crypto/cipher"

	"go.x2ox.com/sorbifolia/cryptutils/cipher"
)

func CCM(block cc.Block, nonceSize, tagSize int) (AEAD, error) {
	return cipher.NewCCMWithNonceAndTagSizes(block, nonceSize, tagSize)
}
