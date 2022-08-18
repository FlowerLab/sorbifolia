package cryptutils

import (
	"crypto/cipher"
)

func CFB(block cipher.Block, iv []byte) CryptStream {
	return streamBlock{
		enc: cipher.NewCFBEncrypter(block, iv),
		dec: cipher.NewCFBDecrypter(block, iv),
	}
}
