package render

import (
	"io"

	"go.x2ox.com/sorbifolia/pyrokinesis"
)

var octetStreamContentType = []byte("application/octet-stream")

func Bytes(data, contentType []byte) Render {
	r := renderPool.Get().(*render)
	r.length = int64(len(data))
	var idx int64

	r.read = func(p []byte) (int, error) {
		if idx == r.length {
			return 0, io.EOF
		}
		n := copy(p, data[idx:r.length])
		idx += int64(n)
		return n, nil
	}
	if contentType == nil {
		contentType = octetStreamContentType
	}
	r.contentType = contentType
	return r
}

var textContentType = []byte("text/plain; charset=utf-8")

func Text[T string | []byte](t T) Render {
	var data []byte
	switch t := any(t).(type) {
	case string:
		data = pyrokinesis.String.ToBytes(t)
	case []byte:
		data = t
	}

	r := renderPool.Get().(*render)
	r.length = int64(len(data))
	var idx int64

	r.read = func(p []byte) (int, error) {
		if idx == r.length {
			return 0, io.EOF
		}
		n := copy(p, data[idx:r.length])
		idx += int64(n)
		return n, nil
	}
	r.contentType = textContentType
	return r
}
