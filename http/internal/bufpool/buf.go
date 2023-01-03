package bufpool

import (
	"bytes"
	"errors"
	"io"
)

type Buffer struct {
	B []byte
}

// Len returns the size of the byte buffer.
func (b *Buffer) Len() int {
	return len(b.B)
}

// ReadFrom implements io.ReaderFrom.
//
// The function appends all the data read from r to b.
func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	p := b.B
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
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

// WriteTo implements io.WriterTo.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.B)
	return int64(n), err
}

// Bytes returns b.B, i.e. all the bytes accumulated in the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
func (b *Buffer) Bytes() []byte {
	return b.B
}

// Write implements io.Writer - it appends p to Buffer.B
func (b *Buffer) Write(p []byte) (int, error) {
	b.B = append(b.B, p...)
	return len(p), nil
}

// WriteByte appends the byte c to the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
//
// The function always returns nil.
func (b *Buffer) WriteByte(c byte) error {
	b.B = append(b.B, c)
	return nil
}

// WriteString appends s to Buffer.B.
func (b *Buffer) WriteString(s string) (int, error) {
	b.B = append(b.B, s...)
	return len(s), nil
}

// Set sets Buffer.B to p.
func (b *Buffer) Set(p []byte) {
	b.B = append(b.B[:0], p...)
}

// SetString sets Buffer.B to s.
func (b *Buffer) SetString(s string) {
	b.B = append(b.B[:0], s...)
}

// String returns string representation of Buffer.B.
func (b *Buffer) String() string {
	return string(b.B)
}

// Reset makes Buffer.B empty.
func (b *Buffer) Reset() {
	b.B = b.B[:0]
}

var ErrWriteLimitExceeded = errors.New("write limit exceeded")

func (b *Buffer) WriteLimit(p []byte, limit int) int {
	i := limit - b.Len()
	switch {
	case i < 0:
		panic("should not mix writes")
	case i == 0:
		return -1
	case i > len(p):
		i = len(p)
	}

	b.B = append(b.B, p[:i]...)
	return i
}

func (b *Buffer) Index(sep []byte) int   { return bytes.Index(b.B, sep) }
func (b *Buffer) Discard(start, end int) { b.B = append(b.B[:start], b.B[end:]...) }
