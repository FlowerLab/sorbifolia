package http

import (
	"arena"
	"io"

	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

type Request struct {
	ver    version.Version
	Method method.Method
	Header RequestHeader
	Body   io.ReadCloser
}

func NewFormData(a *arena.Arena, r Request) {
	fd := arena.New[FormData](a)
	// boundary := r.Header.ContentType
	// fd.Boundary = arena.MakeSlice[byte](a, 2+len(boundary), 2+len(boundary))
	// fd.Boundary[0], fd.Boundary[1] = '-', '-'
	// copy(fd.Boundary[2:], boundary)
	fd.Boundary = r.Header.ContentType // todo: get boundary

	/*
		--Boundary\r\n
			Header\r\n
			Header\r\n
			\r\n
			content\r\n
		--Boundary\r\n
			Header\r\n
			Header\r\n
			\r\n
			content\r\n
		--Boundary--
	*/
}

// FormData multipart/form-data
type FormData struct {
	Boundary []byte

	KV   KVs
	File []io.Closer
}

type FileHeader struct {
	Filename string
	Header   KVs
	Size     int64

	content []byte
	tmpfile string
}
