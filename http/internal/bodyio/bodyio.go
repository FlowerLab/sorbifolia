package bodyio

import (
	"io"
)

var nrc = &null{}

type null struct{}

func (b *null) Read(_ []byte) (int, error) { return 0, io.EOF }

func (b *null) Close() error { return nil }

func Null() io.ReadCloser { return nrc }
