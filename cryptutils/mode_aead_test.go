package cryptutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"testing"
)

var (
	_aes16key = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	_aes24key = append(_aes16key, []byte{17, 18, 19, 20, 21, 22, 23, 24}...)
	_aes32key = append(_aes24key, []byte{25, 26, 27, 28, 29, 30, 31, 32}...)

	_aesKey = [][]byte{_aes16key, _aes24key, _aes32key}
	_iv16   = _aes16key
)

func TestCCM(t *testing.T) {
	t.Parallel()

	testAEAD(t, func(block cipher.Block) AEAD {
		aead, err := CCM(block, 13, 16)
		if err != nil {
			t.Error(err)
		}
		return aead
	})
}

func TestEAX(t *testing.T) {
	t.Parallel()

	testAEAD(t, func(block cipher.Block) AEAD {
		aead, err := EAX(block)
		if err != nil {
			t.Error(err)
		}
		return aead
	})
}

func TestGCM(t *testing.T) {
	t.Parallel()

	testAEAD(t, func(block cipher.Block) AEAD {
		aead, err := GCM(block)
		if err != nil {
			t.Error(err)
		}
		return aead
	})
}

func TestOCB(t *testing.T) {
	t.Parallel()

	testAEAD(t, func(block cipher.Block) AEAD {
		aead, err := OCB(block)
		if err != nil {
			t.Error(err)
		}
		return aead
	})
}

func testAEAD(t *testing.T, fn func(block cipher.Block) AEAD) {
	for _, v := range _aesKey {
		block, err := aes.NewCipher(v)
		if err != nil {
			t.Error(err)
		}
		aead := fn(block)
		nonce := make([]byte, aead.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err)
		}
		plaintext := []byte("exampleplaintext")
		add := []byte("go.x2ox.com/sorbifolia/cryptutils")

		encData := aead.Seal(nil, nonce, plaintext, add)
		decData, err := aead.Open(nil, nonce, encData, add)
		if err != nil {
			t.Error("err")
		}

		if !bytes.Equal(plaintext, decData) {
			t.Error("TestCCM err")
		}
	}
}
