package http

import (
	"io"
)

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on.
type Render interface {
	Render() io.Reader
	Length() int
	ContentType() []byte
}
