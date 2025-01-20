package buffer

import (
	"fmt"
	"io"
)

type Byte struct {
	B []byte
}

func (b *Byte) ReadFrom(r io.Reader) (int64, error) {
	var (
		p      = b.B
		nStart = int64(len(p))
		nMax   = int64(cap(p))
		n      = nStart
	)

	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}

	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
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

func (b *Byte) Write(p []byte) (int, error)        { b.B = append(b.B, p...); return len(p), nil }
func (b *Byte) WriteTo(w io.Writer) (int64, error) { n, err := w.Write(b.B); return int64(n), err }
func (b *Byte) WriteByte(c byte) error             { b.B = append(b.B, c); return nil }
func (b *Byte) WriteString(s string) (int, error)  { b.B = append(b.B, s...); return len(s), nil }

func (b *Byte) Set(p []byte)       { b.B = append(b.B[:0], p...) }
func (b *Byte) SetString(s string) { b.B = append(b.B[:0], s...) }

func (b *Byte) String() string { return string(b.B) }
func (b *Byte) Len() int       { return len(b.B) }
func (b *Byte) Cap() int       { return cap(b.B) }
func (b *Byte) Bytes() []byte  { return b.B }
func (b *Byte) Reset()         { b.B = b.B[:0] }

func (b *Byte) Reader() Reader { return &ByteReader{Byte: b} }

var (
	_ io.Writer       = (*Byte)(nil)
	_ io.StringWriter = (*Byte)(nil)
	_ io.ByteWriter   = (*Byte)(nil)
	_ io.ReaderFrom   = (*Byte)(nil)
	_ fmt.Stringer    = (*Byte)(nil)
	_ Buffer          = (*Byte)(nil)
)
