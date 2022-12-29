package http

import (
	"arena"
	"bytes"
	"io"

	"github.com/pkg/errors"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

type Request struct {
	a      *arena.Arena
	ver    version.Version
	Method method.Method
	Header RequestHeader
	Body   io.ReadCloser
}

func NewFormData(a *arena.Arena, r Request) (*FormData, error) {
	fd := arena.New[FormData](a)
	fd.Boundary = r.Header.ContentType.Boundary() // todo: get boundary
	boundaryLen := len(fd.Boundary)
	if boundaryLen < 1 || boundaryLen > 70 {
		return nil, errors.New("boundary length is not in 1 <= size <= 70")
	}

	length := int(r.Header.ContentLength.Length())
	buf := arena.MakeSlice[byte](a, length, length)

	n, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if n != length {
		return nil, errors.New("content length mismatch")
	}

	for {
		// --boundary
		if !bytes.Equal(buf, fd.Boundary) {
			return nil, errors.New("find not found boundary")
		}
		buf = buf[boundaryLen:]

		// \r\n
		if !bytes.Equal(buf, char.CRLF) {
			return nil, errors.New("find not found \r\n")
		}
		buf = buf[2:]

		ks := arena.MakeSlice[KV](a, 0, 2)

		// header: val
		for {
			i := bytes.Index(buf, char.CRLF)
			if i < 0 {
				return nil, err // has issues
			}
			if i == 0 {
				buf = buf[i+2:]
				break
			}

			kv := arena.New[KV](a)
			kv.ParseHeader(buf[:i])
			ks = append(ks, *kv)

			buf = buf[i+2:]
		}
		kvs := KVs(ks)
		cd := kvs.Get(char.ContentDisposition) // if it has filename, it's a file
		_ = cd
		_ = kvs.Get(char.ContentType) // if it has, it's a file

		// cd.QualityValues()

		// content

	}

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
