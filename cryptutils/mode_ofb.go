package cryptutils

import (
	"crypto/cipher"
)

func OFB(block cipher.Block, iv []byte) CryptStream {
	return streamBlock{enc: cipher.NewOFB(block, iv), dec: cipher.NewOFB(block, iv)}
}
