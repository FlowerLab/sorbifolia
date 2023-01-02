package bodyio

import (
	"io"
)

type block struct {
	preRead []byte
	r       io.Reader
	p       int
}

func (b block) Read(p []byte) (n int, err error) {
	if b.p < len(b.preRead) {
		n = copy(p, b.preRead[b.p:])
		b.p += n
		if b.p < len(b.preRead) {
			return n, nil
		}
	}

	if b.r == nil {
		return 0, io.EOF
	}

	var rn int
	rn, err = b.r.Read(p)
	rn += n
	return
}

func (b block) Close() error { return nil }

func Block(preRead []byte, r io.Reader, length int64) (io.ReadCloser, error) {
	return &block{
		preRead: preRead,
		r:       io.LimitReader(r, length-int64(len(preRead))),
	}, nil
}

func newBlock(preRead []byte, r io.Reader) io.ReadCloser {
	return &block{
		preRead: preRead,
		r:       r,
	}
}
