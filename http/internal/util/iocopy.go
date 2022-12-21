package util

import (
	"io"
	"sync"
)

var copyBufPool = sync.Pool{New: func() any { return make([]byte, 4096) }}

func Copy(w io.Writer, r io.Reader) (n int64, err error) {
	buf := copyBufPool.Get().([]byte)
	n, err = io.CopyBuffer(w, r, buf)
	copyBufPool.Put(buf)
	return
}
