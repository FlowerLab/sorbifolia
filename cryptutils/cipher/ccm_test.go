package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"
)

func TestCcm(t *testing.T) {
	C4A := make([]byte, 524288/8)
	for i := range C4A {
		C4A[i] = byte(i)
	}

	examples := []struct {
		Key        []byte
		Nonce      []byte
		Data       []byte
		PlainText  []byte
		CipherText []byte
		TagLen     int
	}{
		{ // C.1
			[]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f},
			[]byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16},
			[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			[]byte{0x20, 0x21, 0x22, 0x23},
			[]byte{0x71, 0x62, 0x01, 0x5b, 0x4d, 0xac, 0x25, 0x5d},
			4,
		},
		{ // C.2
			[]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f},
			[]byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17},
			[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
			[]byte{0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f},
			[]byte{0xd2, 0xa1, 0xf0, 0xe0, 0x51, 0xea, 0x5f, 0x62, 0x08, 0x1a, 0x77, 0x92, 0x07, 0x3d, 0x59, 0x3d, 0x1f, 0xc6, 0x4f, 0xbf, 0xac, 0xcd},
			6,
		},
		{ // C.3
			[]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f},
			[]byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b},
			[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13},
			[]byte{0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37},

			[]byte{0xe3, 0xb2, 0x01, 0xa9, 0xf5, 0xb7, 0x1a, 0x7a, 0x9b, 0x1c, 0xea, 0xec, 0xcd, 0x97, 0xe7, 0x0b, 0x61, 0x76, 0xaa, 0xd9, 0xa4, 0x42, 0x8a, 0xa5, 0x48, 0x43, 0x92, 0xfb, 0xc1, 0xb0, 0x99, 0x51},
			8,
		},
		{ // C.4
			[]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f},
			[]byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c},
			C4A,
			[]byte{0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f},

			[]byte{0x69, 0x91, 0x5d, 0xad, 0x1e, 0x84, 0xc6, 0x37, 0x6a, 0x68, 0xc2, 0x96, 0x7e, 0x4d, 0xab, 0x61, 0x5a, 0xe0, 0xfd, 0x1f, 0xae, 0xc4, 0x4c, 0xc4, 0x84, 0x82, 0x85, 0x29, 0x46, 0x3c, 0xcf, 0x72, 0xb4, 0xac, 0x6b, 0xec, 0x93, 0xe8, 0x59, 0x8e, 0x7f, 0x0d, 0xad, 0xbc, 0xea, 0x5b},
			14,
		},
	}

	for _, v := range examples {
		c, err := aes.NewCipher(v.Key)
		if err != nil {
			t.Fatal(err)
		}

		aead, err := NewCCMWithNonceAndTagSizes(c, len(v.Nonce), v.TagLen)
		if err != nil {
			t.Fatal(err)
		}

		CipherText := aead.Seal(nil, v.Nonce, v.PlainText, v.Data)

		if !bytes.Equal(v.CipherText, CipherText) {
			t.Fatal("fail")
		}

		PlainText, err := aead.Open(nil, v.Nonce, v.CipherText, v.Data)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(v.PlainText, PlainText) {
			t.Fatal("fail")
		}
	}
}

func TestNewCCMWithNonceAndTagSizes(t *testing.T) {
	block, err := aes.NewCipher([]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("", func(t *testing.T) {
		if _, err = NewCCMWithNonceAndTagSizes(block, 6, 4); err == nil {
			t.Error("expected")
		}
	})
	t.Run("", func(t *testing.T) {
		if _, err = NewCCMWithNonceAndTagSizes(block, 7, 3); err == nil {
			t.Error("expected")
		}
	})
	t.Run("", func(t *testing.T) {
		if _, err = NewCCMWithNonceAndTagSizes(testBlock(12), 6, 4); err == nil {
			t.Error("expected")
		}
	})
	t.Run("", func(t *testing.T) {
		a, _ := NewCCMWithNonceAndTagSizes(block, 7, 4)
		if a.NonceSize() != 7 || a.Overhead() != 4 {
			t.Error("expected")
		}

		defer func() { _ = recover() }()
		a.Seal(nil, nil, nil, nil)
		t.Error("expected")
	})
	t.Run("", func(t *testing.T) {
		aead, _ := NewCCMWithNonceAndTagSizes(block, 7, 4)
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
	})
	t.Run("", func(t *testing.T) {
		a, _ := NewCCMWithNonceAndTagSizes(block, 7, 4)
		if a.NonceSize() != 7 || a.Overhead() != 4 {
			t.Error("expected")
		}

		defer func() { _ = recover() }()
		if _, err := a.Open(nil, nil, nil, nil); err == nil {
			t.Error("expected")
		}
	})
}

type testBlock int

func (t testBlock) BlockSize() int          { return int(t) }
func (t testBlock) Encrypt(dst, src []byte) {}
func (t testBlock) Decrypt(dst, src []byte) {}

func TestCCM(t *testing.T) {
	var testDataRfc3610 = []struct {
		key        string
		nonce      string
		adata      string
		plaintext  string
		ciphertext string
	}{
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "00000003020100a0a1a2a3a4a5", adata: "0001020304050607", plaintext: "08090a0b0c0d0e0f101112131415161718191a1b1c1d1e", ciphertext: "588c979a61c663d2f066d0c2c0f989806d5f6b61dac38417e8d12cfdf926e0"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "00000004030201a0a1a2a3a4a5", adata: "0001020304050607", plaintext: "08090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f", ciphertext: "72c91a36e135f8cf291ca894085c87e3cc15c439c9e43a3ba091d56e10400916"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "00000005040302a0a1a2a3a4a5", adata: "0001020304050607", plaintext: "08090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20", ciphertext: "51b1e5f44a197d1da46b0f8e2d282ae871e838bb64da8596574adaa76fbd9fb0c5"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "00000006050403a0a1a2a3a4a5", adata: "000102030405060708090a0b", plaintext: "0c0d0e0f101112131415161718191a1b1c1d1e", ciphertext: "a28c6865939a9a79faaa5c4c2a9d4a91cdac8c96c861b9c9e61ef1"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "00000007060504a0a1a2a3a4a5", adata: "000102030405060708090a0b", plaintext: "0c0d0e0f101112131415161718191a1b1c1d1e1f", ciphertext: "dcf1fb7b5d9e23fb9d4e131253658ad86ebdca3e51e83f077d9c2d93"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "00000008070605a0a1a2a3a4a5", adata: "000102030405060708090a0b", plaintext: "0c0d0e0f101112131415161718191a1b1c1d1e1f20", ciphertext: "6fc1b011f006568b5171a42d953d469b2570a4bd87405a0443ac91cb94"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "00000009080706a0a1a2a3a4a5", adata: "0001020304050607", plaintext: "08090a0b0c0d0e0f101112131415161718191a1b1c1d1e", ciphertext: "0135d1b2c95f41d5d1d4fec185d166b8094e999dfed96c048c56602c97acbb7490"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "0000000a090807a0a1a2a3a4a5", adata: "0001020304050607", plaintext: "08090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f", ciphertext: "7b75399ac0831dd2f0bbd75879a2fd8f6cae6b6cd9b7db24c17b4433f434963f34b4"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "0000000b0a0908a0a1a2a3a4a5", adata: "0001020304050607", plaintext: "08090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20", ciphertext: "82531a60cc24945a4b8279181ab5c84df21ce7f9b73f42e197ea9c07e56b5eb17e5f4e"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "0000000c0b0a09a0a1a2a3a4a5", adata: "000102030405060708090a0b", plaintext: "0c0d0e0f101112131415161718191a1b1c1d1e", ciphertext: "07342594157785152b074098330abb141b947b566aa9406b4d999988dd"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "0000000d0c0b0aa0a1a2a3a4a5", adata: "000102030405060708090a0b", plaintext: "0c0d0e0f101112131415161718191a1b1c1d1e1f", ciphertext: "676bb20380b0e301e8ab79590a396da78b834934f53aa2e9107a8b6c022c"},
		{key: "c0c1c2c3c4c5c6c7c8c9cacbcccdcecf", nonce: "0000000e0d0c0ba0a1a2a3a4a5", adata: "000102030405060708090a0b", plaintext: "0c0d0e0f101112131415161718191a1b1c1d1e1f20", ciphertext: "c0ffa0d6f05bdb67f24d43a4338d2aa4bed7b20e43cd1aa31662e7ad65d6db"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "00412b4ea9cdbe3c9696766cfa", adata: "0be1a88bace018b1", plaintext: "08e8cf97d820ea258460e96ad9cf5289054d895ceac47c", ciphertext: "4cb97f86a2a4689a877947ab8091ef5386a6ffbdd080f8e78cf7cb0cddd7b3"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "0033568ef7b2633c9696766cfa", adata: "63018f76dc8a1bcb", plaintext: "9020ea6f91bdd85afa0039ba4baff9bfb79c7028949cd0ec", ciphertext: "4ccb1e7ca981befaa0726c55d378061298c85c92814abc33c52ee81d7d77c08a"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "00103fe41336713c9696766cfa", adata: "aa6cfa36cae86b40", plaintext: "b916e0eacc1c00d7dcec68ec0b3bbb1a02de8a2d1aa346132e", ciphertext: "b1d23a2220ddc0ac900d9aa03c61fcf4a559a4417767089708a776796edb723506"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "00764c63b8058e3c9696766cfa", adata: "d0d0735c531e1becf049c244", plaintext: "12daac5630efa5396f770ce1a66b21f7b2101c", ciphertext: "14d253c3967b70609b7cbb7c499160283245269a6f49975bcadeaf"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "00f8b678094e3b3c9696766cfa", adata: "77b60f011c03e1525899bcae", plaintext: "e88b6a46c78d63e52eb8c546efb5de6f75e9cc0d", ciphertext: "5545ff1a085ee2efbf52b2e04bee1e2336c73e3f762c0c7744fe7e3c"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "00d560912d3f703c9696766cfa", adata: "cd9044d2b71fdb8120ea60c0", plaintext: "6435acbafb11a82e2f071d7ca4a5ebd93a803ba87f", ciphertext: "009769ecabdf48625594c59251e6035722675e04c847099e5ae0704551"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "0042fff8f1951c3c9696766cfa", adata: "d85bc7e69f944fb8", plaintext: "8a19b950bcf71a018e5e6701c91787659809d67dbedd18", ciphertext: "bc218daa947427b6db386a99ac1aef23ade0b52939cb6a637cf9bec2408897c6ba"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "00920f40e56cdc3c9696766cfa", adata: "74a0ebc9069f5b37", plaintext: "1761433c37c5a35fc1f39f406302eb907c6163be38c98437", ciphertext: "5810e6fd25874022e80361a478e3e9cf484ab04f447efff6f0a477cc2fc9bf548944"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "0027ca0c7120bc3c9696766cfa", adata: "44a3aa3aae6475ca", plaintext: "a434a8e58500c6e41530538862d686ea9e81301b5ae4226bfa", ciphertext: "f2beed7bc5098e83feb5b31608f8e29c38819a89c8e776f1544d4151a4ed3a8b87b9ce"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "005b8ccbcd9af83c9696766cfa", adata: "ec46bb63b02520c33c49fd70", plaintext: "b96b49e21d621741632875db7f6c9243d2d7c2", ciphertext: "31d750a09da3ed7fddd49a2032aabf17ec8ebf7d22c8088c666be5c197"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "003ebe94044b9a3c9696766cfa", adata: "47a65ac78b3d594227e85e71", plaintext: "e2fcfbb880442c731bf95167c8ffd7895e337076", ciphertext: "e882f1dbd38ce3eda7c23f04dd65071eb41342acdf7e00dccec7ae52987d"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "003ebe94044b9a3c9696766cfa", adata: "47a65ac78b3d594227e85e71", plaintext: "e2fcfbb880442c731bf95167c8ffd7895e337076", ciphertext: "e882f1dbd38ce3eda7c23f04dd65071eb41342acdf7e00dccec7ae52987d"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "003ebe94044b9a3c9696766cfa", adata: "47a65ac78b3d594227e85e71", plaintext: "e2fcfbb880442c731bf95167c8ffd7895e337076", ciphertext: "e882f1dbd38ce3eda7c23f04dd65071eb41342acdf7e00dccec7ae52987d"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "008d493b30ae8b3c9696766cfa", adata: "6e37a6ef546d955d34ab6059", plaintext: "abf21c0b02feb88f856df4a37381bce3cc128517d4", ciphertext: "f32905b88a641b04b9c9ffb58cc390900f3da12ab16dce9e82efa16da62059"},
		{key: "d7828d13b2b0bdc325a76236df93cc6b", nonce: "008d493b30ae8b3c9696766cfa", adata: "6e37a6ef546d955d34ab6059", plaintext: "abf21c0b02feb88f856df4a37381bce3cc128517d4", ciphertext: "f32905b88a641b04b9c9ffb58cc390900f3da12ab16dce9e82efa16da62059"},
	}

	decodeAndCheck := func(vv string, i int) (rv []byte) {
		var err error
		rv, err = hex.DecodeString(vv)
		if err != nil {
			t.Errorf("AesCCM FATAL ERROR: Unable to setup AES, input hex failed to parse, %v test:#%d", vv, i)
			return
		}
		return
	}

	for i, v := range testDataRfc3610 {
		key := decodeAndCheck(v.key, i)
		nonce := decodeAndCheck(v.nonce, i)
		adata := decodeAndCheck(v.adata, i)
		plaintext := decodeAndCheck(v.plaintext, i)
		Aes, err := aes.NewCipher(key)
		if err != nil {
			t.Errorf("AesCCM FATAL ERROR: Unable to setup AES with given key")
			return
		}

		TagLength := hex.DecodedLen(len(v.ciphertext)) - len(plaintext)
		AesCCM, err := NewCCMWithNonceAndTagSizes(Aes, len(nonce), TagLength)
		if err != nil {
			t.Fatal(err)
		}

		ct := AesCCM.Seal(nil, nonce, plaintext, adata)
		tmp := fmt.Sprintf("%x", ct)
		if !strings.EqualFold(v.ciphertext, tmp) {
			t.Errorf("AesCCM Test #%d: got\t%s, expected\t%s", i, tmp, v.ciphertext)
			continue
		}

		plaintext2, err := AesCCM.Open(nil, nonce, ct, adata)
		if err != nil {
			t.Errorf("AesCCM Test #%d: Open failed when it should have succeded: %v", i, err)
			continue
		}

		if !bytes.Equal(plaintext, plaintext2) {
			t.Errorf("AesCCM Test #%d: got %x expected %x, failed to properly recover original data", i, plaintext2, plaintext)
			continue
		}

		for j := 0; j < 8; j++ {
			onebit := byte(1 << uint(j))
			for pos := 0; pos < len(nonce) || pos < len(ct); pos++ {
				if pos < len(nonce) {
					nonce[pos] ^= onebit
					if _, err := AesCCM.Open(nil, nonce, ct, adata); err == nil {
						t.Errorf("AesCCM Test #%d: Altered nonce, should have failed open, pos=%d j=%d", i, pos, j)
					}
					nonce[pos] ^= onebit
				}

				if pos < len(ct) {
					ct[pos] ^= onebit
					if _, err := AesCCM.Open(nil, nonce, ct, adata); err == nil {
						t.Errorf("AesCCM Test #%d: Altered ct, should have failed open, pos=%d j=%d", i, pos, j)
					}
					ct[pos] ^= onebit
				}

				if pos < len(adata) && len(adata) > 0 {
					adata[pos] ^= onebit
					if _, err := AesCCM.Open(nil, nonce, ct, adata); err == nil {
						t.Errorf("AesCCM Test #%d: Alterd adata, should have failed open, pos=%d j=%d", i, pos, j)
					}
					adata[pos] ^= onebit
				}

			}
		}
	}
}

func TestCcmGetTag(t *testing.T) {
	t.Run("", func(t *testing.T) {
		block, err := aes.NewCipher([]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f})
		if err != nil {
			t.Fatal(err)
		}

		a, _ := NewCCMWithNonceAndTagSizes(block, 7, 4)
		if a.NonceSize() != 7 || a.Overhead() != 4 {
			t.Error("expected")
		}
		_ccm, ok := a.(*ccm)
		if !ok {
			t.Error("expected")
		}

		// Formatting of the Counter Blocks are defined in A.3.
		Ctr := make([]byte, 16)                    // Ctr0
		Ctr[0] = byte(15 - _ccm.nonceSize - 1)     // [q-1]3
		copy(Ctr[1:], []byte{1, 2, 3, 4, 5, 6, 7}) // N

		S0 := make([]byte, 16) // S0
		_ccm.c.Encrypt(S0, Ctr)

		Ctr[15] = 1 // Ctr1

		if runtime.GOOS != "windows" {
			// GitHub Action Windows ThreadSanitizer failed to allocate
			return
		}

		add := make([]byte, 1<<31)
		plaintext := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

		data := _ccm.getTag(Ctr, add, plaintext)
		fmt.Println(data)
	})

	t.Run("", func(t *testing.T) {
		block, err := aes.NewCipher([]byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f})
		if err != nil {
			t.Fatal(err)
		}

		a, _ := NewCCMWithNonceAndTagSizes(block, 7, 4)
		if a.NonceSize() != 7 || a.Overhead() != 4 {
			t.Error("expected")
		}
		_ccm, ok := a.(*ccm)
		if !ok {
			t.Error("expected")
		}

		// Formatting of the Counter Blocks are defined in A.3.
		Ctr := make([]byte, 16)                    // Ctr0
		Ctr[0] = byte(15 - _ccm.nonceSize - 1)     // [q-1]3
		copy(Ctr[1:], []byte{1, 2, 3, 4, 5, 6, 7}) // N

		S0 := make([]byte, 16) // S0
		_ccm.c.Encrypt(S0, Ctr)

		Ctr[15] = 1 // Ctr1

		plaintext := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

		_ccm.getTag(Ctr, nil, plaintext)
	})
}
