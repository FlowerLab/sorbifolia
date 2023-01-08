package httpmessage

import (
	"errors"
	"io"
)

func (r *Request) Read(p []byte) (n int, err error) {
	if !r.state.Readable() {
		return 0, errors.New("? TODO: here we need to think about how to deal with")
	}
	if length := r.buf.Len(); length == 0 || length == r.rp {
		return 0, io.EOF
	}
	n = copy(p, r.buf.B[r.rp:])
	r.rp += n
	return
}
