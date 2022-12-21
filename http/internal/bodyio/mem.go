//go:build goexperiment.arenas

package bodyio

import (
	"arena"
	"errors"
	"io"
)

type mem struct {
	buf []byte
	p   int
}

func (m *mem) Read(p []byte) (n int, err error) {
	if m.p == len(m.buf) {
		return 0, io.EOF
	}
	n = copy(p, m.buf[m.p:])
	m.p += n
	return
}

func (m *mem) Close() error { return nil }

func Memory(a *arena.Arena, preRead []byte, r io.Reader, length int64) (io.ReadCloser, error) {
	bf := arena.New[mem](a)
	bf.buf = arena.MakeSlice[byte](a, int(length), int(length))
	cn := copy(bf.buf, preRead)
	r = io.LimitReader(r, length-int64(len(preRead)))

	for int64(cn) < length {
		n, err := r.Read(bf.buf[cn:])
		cn += n

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	if int64(cn) != length {
		return nil, errors.New("length mismatch")
	}

	return bf, nil
}
