package bufpool

import (
	"io"
)

type ReadBuffer struct {
	Buffer
	P int
}

func (r *ReadBuffer) Read(p []byte) (n int, err error) {
	if r.P == r.Len() {
		return 0, io.EOF
	}
	n = copy(p, r.B[r.P:])
	r.P += n
	return
}

func (r *ReadBuffer) Reset() {
	r.B = r.B[:0]
	r.P = 0
}

func (r *ReadBuffer) Release() {
	r.Reset()
	readBufPool.Put(r)
}
