package render

import (
	"arena"
	"encoding/json"
	"io"

	"go.x2ox.com/sorbifolia/http"
	"go.x2ox.com/sorbifolia/http/internal/util"
)

var jsonContentType = []byte("application/json; charset=utf-8")

type render struct {
	read        func(p []byte) (n int, err error)
	close       func() error
	contentType []byte
	buf         *util.Buffer
	err         error
}

func JSON(a *arena.Arena, data any) http.Render {
	r := arena.New[render](a)
	r.buf = arena.New[util.Buffer](a)
	r.buf.A = a
	r.err = json.NewEncoder(r.buf).Encode(data)

	r.read = func(p []byte) (n int, err error) {
		if r.err != nil {
			return 0, err
		}
		return r.buf.Read(p)
	}
	r.close = func() error { return nil }
	r.contentType = jsonContentType
	return r
}

func (r *render) Render() io.Reader                { return r }
func (r *render) ContentType() []byte              { return r.contentType }
func (r *render) Length() int                      { return r.buf.Len() }
func (r *render) Read(p []byte) (n int, err error) { return r.read(p) }
func (r *render) Close() error                     { return r.close() }

var (
	_ io.ReadCloser = (*render)(nil)
	_ http.Render   = (*render)(nil)
)
