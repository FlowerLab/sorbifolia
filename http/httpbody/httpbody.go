package httpbody

import (
	"io"
)

type HTTPBody interface {
	BodyReader() io.ReadCloser
	BodyWriter() io.WriteCloser
}
