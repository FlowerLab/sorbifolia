package util

import (
	"io"
	"sync"
)

var copyBufPool = sync.Pool{New: func() any { return make([]byte, 4096) }}

func Copy(w io.Writer, r io.Reader) (n int64, err error) {
	buf := copyBufPool.Get()
	b := buf.([]byte)
	n, err = io.CopyBuffer(w, r, b)
	copyBufPool.Put(buf)
	return
}

func ReadAll(b []byte, r io.Reader) (int, error) {
	var i int
	for {
		n, err := r.Read(b[i:])
		i += n
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return i, err
		}
	}
}
