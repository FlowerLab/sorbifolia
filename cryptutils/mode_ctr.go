package cryptutils

import (
	"crypto/cipher"
)

func CTR(block cipher.Block, iv []byte) CryptStream {
	// stream := cipher.NewCTR(block, iv)
	// return streamBlock{enc: stream, dec: stream}
	// stream := cipher.NewCTR(block, iv)
	return streamBlock{enc: cipher.NewCTR(block, iv), dec: cipher.NewCTR(block, iv)}
}
