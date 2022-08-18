package cryptutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"testing"
)

func TestOFB(t *testing.T) {
	testAESStream(t, func(block cipher.Block, iv []byte) CryptStream {
		return OFB(block, iv)
	})
}

func TestCFB(t *testing.T) {
	testAESStream(t, func(block cipher.Block, iv []byte) CryptStream {
		return CFB(block, iv)
	})
}

func TestCTR(t *testing.T) {
	testAESStream(t, func(block cipher.Block, iv []byte) CryptStream {
		return CTR(block, iv)
	})
}

func testStream(t *testing.T, cs CryptStream) {
	var (
		src = []byte{1, 2, 3, 4, 5, 6, 7}

		encData = make([]byte, len(src))
		decData = make([]byte, len(src))
	)

	cryptStreamEncrypt(cs, encData, src)
	cryptStreamDecrypt(cs, decData, encData)

	if bytes.Compare(src, decData) != 0 {
		t.Error("testStream err")
	}
}

func testAESStream(t *testing.T, fn func(block cipher.Block, iv []byte) CryptStream) {
	for _, v := range _aesKey {
		block, err := aes.NewCipher(v)
		if err != nil {
			t.Error(err)
		}
		testStream(t, fn(block, _iv16))
	}
}

func cryptStreamEncrypt(cb CryptStream, dst, src []byte) { cb.Encrypt(dst, src) }
func cryptStreamDecrypt(cb CryptStream, dst, src []byte) { cb.Decrypt(dst, src) }
