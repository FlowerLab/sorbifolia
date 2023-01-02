package http

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

type Request struct {
	ver    version.Version
	Method method.Method
	Header RequestHeader
	Body   io.ReadCloser
}

func NewFormData(r Request) (*FormData, error) {
	fd := &FormData{}
	fd.Boundary = r.Header.ContentType.Boundary()
	boundaryLen := len(fd.Boundary)
	if boundaryLen < 1 || boundaryLen > 70 {
		return nil, errors.New("boundary length is not in 1 <= size <= 70")
	}

	length := int(r.Header.ContentLength.Length())
	buf := make([]byte, length, length)

	n, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if n != length {
		return nil, errors.New("content length mismatch")
	}

	for {
		if !bytes.Equal(buf, fd.Boundary) { // --boundary
			if !bytes.Equal(buf, []byte("--")) { // --boundary--
				break
			}
			return nil, errors.New("find not found boundary")
		}
		buf = buf[boundaryLen:]

		if !bytes.Equal(buf, char.CRLF) { // \r\n
			return nil, errors.New("find not found \r\n")
		}
		buf = buf[2:]

		ks := make([]KV, 0, 2)

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

			// kv := arena.New[KV](a)
			kv := &KV{}
			kv.ParseHeader(buf[:i])
			ks = append(ks, *kv)

			buf = buf[i+2:]
		}
		kvs := KVs(ks)
		cd := kvs.Get(char.ContentDisposition)

		hcd := httpheader.ContentDisposition(cd.V)
		filename := hcd.Filename()
		if len(filename) == 0 {
			kv := &KV{}
			// kv := arena.New[KV](a)
			kv.K = hcd.Name()

			i := bytes.Index(buf, fd.Boundary)
			if i < 0 { // --boundary
				return nil, errors.New("find not found boundary")
			}

			kv.V = buf[:i]
			if len(fd.KV)+1 < cap(fd.KV) {
				// arr := arena.MakeSlice[KV](a, len(fd.KV), len(fd.KV)+1)
				arr := make([]KV, len(fd.KV), len(fd.KV)+1)
				copy(arr, fd.KV)
				fd.KV = arr
			}
			fd.KV = append(fd.KV, *kv)

			buf = buf[i+len(fd.Boundary):]
			continue
		}

		i := bytes.Index(buf, fd.Boundary)
		if i < 0 { // --boundary
			return nil, errors.New("find not found boundary")
		}
		fh := FileHeader{
			Name:     filename,
			Filename: hcd.Name(),
			Size:     int64(i),
			Header:   kvs,
			content:  buf[:i],
		}

		fd.File = append(fd.File, fh)
		buf = buf[i+len(fd.Boundary):]
	}
	return fd, nil
}

// FormData multipart/form-data
type FormData struct {
	Boundary []byte

	KV   KVs
	File []FileHeader
}

type FileHeader struct {
	Name     []byte
	Filename []byte
	Size     int64
	Header   KVs

	content []byte
	tmpfile string
}
