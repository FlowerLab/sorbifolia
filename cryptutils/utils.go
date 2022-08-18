package cryptutils

import (
	"crypto/cipher"

	"go.x2ox.com/sorbifolia/cryptutils/padding"
)

type cryptBlock struct {
	block    cipher.Block
	pad      padding.Padding
	enc, dec func([]byte, []byte)
}

func (c cryptBlock) Encrypt(src []byte) []byte {
	data, _ := c.pad.Pad(src, c.block.BlockSize())
	buf := make([]byte, len(data))
	c.enc(buf, data)
	return buf
}

func (c cryptBlock) Decrypt(src []byte) []byte {
	buf := make([]byte, len(src))
	c.dec(buf, src)
	data, _ := c.pad.UnPad(buf, c.block.BlockSize())
	return data
}

type streamBlock struct{ enc, dec cipher.Stream }

func (c streamBlock) Encrypt(dst, src []byte) { c.enc.XORKeyStream(dst, src) }
func (c streamBlock) Decrypt(dst, src []byte) { c.dec.XORKeyStream(dst, src) }
