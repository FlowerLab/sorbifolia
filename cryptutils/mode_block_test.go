package cryptutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"go.x2ox.com/sorbifolia/cryptutils/padding"
)

func TestCBC(t *testing.T) {
	testAESBlock(t, func(block cipher.Block, iv []byte, pad padding.Padding) CryptBlock {
		return CBC(block, iv, pad)
	})
}

func TestECB(t *testing.T) {
	testAESBlock(t, func(block cipher.Block, iv []byte, pad padding.Padding) CryptBlock {
		return ECB(block, pad)
	})
}

func testBlock(t *testing.T, cb CryptBlock) {
	src := []byte{1, 2, 3, 4, 5, 6, 7}
	encData := cryptBlockEncrypt(cb, src)
	decData := cryptBlockDecrypt(cb, encData)

	if bytes.Compare(src, decData) != 0 {
		t.Error("testBlock err")
	}
}

func testAESBlock(t *testing.T, fn func(block cipher.Block, iv []byte, pad padding.Padding) CryptBlock) {
	for _, v := range _aesKey {
		block, err := aes.NewCipher(v)
		if err != nil {
			t.Error(err)
		}

		testBlock(t, fn(block, _iv16, padding.PKCS7{}))
		testBlock(t, fn(block, _iv16, padding.ZeroPadding{}))
		testBlock(t, fn(block, _iv16, padding.ISO10126{}))
		testBlock(t, fn(block, _iv16, padding.ANSIx923{}))
		testBlock(t, fn(block, _iv16, nil))
	}
}

func cryptBlockEncrypt(cb CryptBlock, data []byte) []byte { return cb.Encrypt(data) }
func cryptBlockDecrypt(cb CryptBlock, data []byte) []byte { return cb.Decrypt(data) }
