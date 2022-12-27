package render

import (
	"io"
	"sync"
)

type render struct {
	read        func(p []byte) (n int, err error)
	close       func() error
	contentType []byte
	length      int64
}

func (r *render) Render() io.Reader                { return r }
func (r *render) ContentType() []byte              { return r.contentType }
func (r *render) Length() int64                    { return r.length }
func (r *render) Read(p []byte) (n int, err error) { return r.read(p) }
func (r *render) Reset()                           { r.read = nil; r.close = nil; r.contentType = nil }
func (r *render) Close() (err error) {
	if r.close != nil {
		err = r.close()
	}
	r.Reset()
	renderPool.Put(r)
	return
}

var (
	_ io.ReadCloser = (*render)(nil)
	_ Render        = (*render)(nil)

	renderPool = sync.Pool{New: func() any { return &render{} }}
)

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on.
type Render interface {
	Render() io.Reader
	Length() int64
	ContentType() []byte
}
