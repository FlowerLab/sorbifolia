package cryptutils

import (
	cc "crypto/cipher"

	"go.x2ox.com/sorbifolia/cryptutils/cipher"
	"go.x2ox.com/sorbifolia/cryptutils/padding"
)

func ECB(block cc.Block, pad padding.Padding) CryptBlock {
	if pad == nil {
		pad = padding.PKCS7{}
	}
	enc := cipher.NewECBEncrypter(block)
	dec := cipher.NewECBDecrypter(block)
	return &cryptBlock{block: block, pad: pad, enc: enc.CryptBlocks, dec: dec.CryptBlocks}
}
