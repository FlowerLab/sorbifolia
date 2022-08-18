package cryptutils

import (
	"crypto/cipher"

	"go.x2ox.com/sorbifolia/cryptutils/padding"
)

func CBC(block cipher.Block, iv []byte, pad padding.Padding) CryptBlock {
	if pad == nil {
		pad = padding.PKCS7{}
	}
	enc := cipher.NewCBCEncrypter(block, iv)
	dec := cipher.NewCBCDecrypter(block, iv)
	return &cryptBlock{block: block, pad: pad, enc: enc.CryptBlocks, dec: dec.CryptBlocks}
}
