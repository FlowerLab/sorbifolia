package cbcmac

import (
	"crypto/cipher"
)

type MAC struct {
	ci []byte
	p  int
	c  cipher.Block
}

func New(c cipher.Block) *MAC {
	return &MAC{
		c:  c,
		ci: make([]byte, c.BlockSize()),
	}
}

func (m *MAC) Reset() {
	for i := range m.ci {
		m.ci[i] = 0
	}
	m.p = 0
}

func (m *MAC) Write(p []byte) (n int, err error) {
	for _, c := range p {
		if m.p >= len(m.ci) {
			m.c.Encrypt(m.ci, m.ci)
			m.p = 0
		}
		m.ci[m.p] ^= c
		m.p++
	}
	return len(p), nil
}

// PadZero emulates zero byte padding.
func (m *MAC) PadZero() {
	if m.p != 0 {
		m.c.Encrypt(m.ci, m.ci)
		m.p = 0
	}
}

func (m *MAC) Sum(in []byte) []byte {
	if m.p != 0 {
		m.c.Encrypt(m.ci, m.ci)
		m.p = 0
	}
	return append(in, m.ci...)
}

func (m *MAC) Size() int      { return len(m.ci) }
func (m *MAC) BlockSize() int { return 16 }
