// Package cryptutils has been deprecated since Go 1.24
package cryptutils

import (
	"crypto/cipher"

	cc "go.x2ox.com/sorbifolia/cryptutils/cipher"
	"go.x2ox.com/sorbifolia/cryptutils/padding"
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

func CBC(block cipher.Block, iv []byte, pad padding.Padding) CryptBlock {
	if pad == nil {
		pad = padding.PKCS7{}
	}
	enc := cipher.NewCBCEncrypter(block, iv)
	dec := cipher.NewCBCDecrypter(block, iv)
	return &cryptBlock{block: block, pad: pad, enc: enc.CryptBlocks, dec: dec.CryptBlocks}
}

func CCM(block cipher.Block, nonceSize, tagSize int) (AEAD, error) {
	return cc.NewCCMWithNonceAndTagSizes(block, nonceSize, tagSize)
}
func CFB(block cipher.Block, iv []byte) CryptStream {
	return streamBlock{
		enc: cipher.NewCFBEncrypter(block, iv),
		dec: cipher.NewCFBDecrypter(block, iv),
	}
}
func CTR(block cipher.Block, iv []byte) CryptStream {
	// stream := cipher.NewCTR(block, iv)
	// return streamBlock{enc: stream, dec: stream}
	// stream := cipher.NewCTR(block, iv)
	return streamBlock{enc: cipher.NewCTR(block, iv), dec: cipher.NewCTR(block, iv)}
}

func EAX(block cipher.Block) (AEAD, error) {
	return cc.NewEAX(block)
}

func ECB(block cipher.Block, pad padding.Padding) CryptBlock {
	if pad == nil {
		pad = padding.PKCS7{}
	}
	enc := cc.NewECBEncrypter(block)
	dec := cc.NewECBDecrypter(block)
	return &cryptBlock{block: block, pad: pad, enc: enc.CryptBlocks, dec: dec.CryptBlocks}
}

func GCM(block cipher.Block) (AEAD, error) {
	return cipher.NewGCM(block)
}

func OCB(block cipher.Block) (AEAD, error) {
	return cc.NewOCB(block)
}

func OFB(block cipher.Block, iv []byte) CryptStream {
	return streamBlock{enc: cipher.NewOFB(block, iv), dec: cipher.NewOFB(block, iv)}
}
