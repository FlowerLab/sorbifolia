package httpbody

import (
	"io"
)

var (
	_ io.ReadWriteCloser = (*Chunked)(nil)
	_ io.ReadWriteCloser = (*Memory)(nil)
	_ io.ReadWriteCloser = (*TempFile)(nil)
	_ io.ReadWriteCloser = (*nobody)(nil)
)
