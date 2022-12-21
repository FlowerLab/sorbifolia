//go:build goexperiment.arenas

package bodyio

import (
	"arena"
	"io"
	"os"
)

type fs struct {
	file     *os.File
	filename string
	close    bool
}

func (b *fs) Read(p []byte) (n int, err error) {
	if b.close {
		return 0, io.EOF
	}
	if b.file == nil {
		if b.file, err = os.Open(b.filename); err != nil {
			return 0, err
		}
	}
	return b.file.Read(p)
}

func (b *fs) Close() error {
	if b.close {
		return nil
	}
	b.close = true
	if err := b.file.Close(); err != nil {
		return err
	}
	return os.Remove(b.filename)
}

func File(a *arena.Arena, preRead []byte, r io.Reader, length int64) (io.ReadCloser, error) {
	file, err := os.CreateTemp("", "http-request-*")
	if err != nil {
		return nil, err
	}
	var n int
	if len(preRead) != 0 {
		if n, err = file.Write(preRead); err != nil {
			return nil, err
		}
	}
	r = io.LimitReader(r, length-int64(len(preRead)))

	var cn int64
	if cn, err = io.Copy(file, r); err != nil {
		return nil, err
	}
	if cn+int64(n) != length {
		return nil, err
	}

	if err = file.Sync(); err != nil {
		return nil, err
	}

	bf := arena.New[fs](a)
	bf.filename = file.Name()

	return bf, nil
}
