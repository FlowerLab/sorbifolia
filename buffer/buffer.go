package buffer

import (
	"io"
)

type Buffer interface {
	Bytes() []byte
	Len() int
	Cap() int
	Reset()

	Reader() Reader
}

type Reader interface {
	Buffer
	io.ReadSeeker
}
