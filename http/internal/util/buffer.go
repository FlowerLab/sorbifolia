//go:build goexperiment.arenas

package util

import (
	"arena"
	"io"

	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type Buffer struct {
	A    *arena.Arena
	B    []byte
	r, w int
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.r < len(b.B) {
		n = copy(p, b.B[b.r:])
		b.r += n
		return n, nil
	}
	return 0, io.EOF
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	if l := len(p); cap(b.B)-b.w < l {
		buf := b.makeSlice(l, l)
		copy(buf, b.B[:b.w])
		b.B = buf
	}
	return copy(b.B[b.w:], p), nil
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.B)
	return int64(n), err
}

func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	p := b.B
	nStart := int64(len(p))
	nMax := cap(p)
	n = nStart
	if nMax == 0 {
		nMax = 64
		p = b.makeSlice(nMax, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == int64(nMax) {
			nMax *= 2
			bNew := b.makeSlice(nMax, nMax)
			copy(bNew, p)
			p = bNew
		}
		var nn int
		nn, err = r.Read(p[n:])
		n += int64(nn)
		if err != nil {
			b.B = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

func (b *Buffer) Len() int       { return len(b.B) }
func (b *Buffer) Bytes() []byte  { return b.B }
func (b *Buffer) String() string { return pyrokinesis.Bytes.ToString(b.B) }
func (b *Buffer) Reset()         { b.B, b.r, b.w = b.B[:0], 0, 0 }

func (b *Buffer) makeSlice(length, capacity int) []byte {
	if b.A == nil {
		return make([]byte, length, capacity)
	}
	return arena.MakeSlice[byte](b.A, length, capacity)
}
