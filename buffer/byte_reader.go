package buffer

import (
	"errors"
	"io"
	"io/fs"
)

type ByteReader struct {
	*Byte
	offset int64
}

func (r *ByteReader) Close() error {
	r.Byte.Free()
	r.B = nil
	r.offset = 0
	return nil
}

func (r *ByteReader) Read(p []byte) (n int, err error) {
	if r.Byte == nil {
		return 0, fs.ErrClosed
	}

	if n = copy(p, r.B[r.offset:]); n < len(p) {
		err = io.EOF
	}
	r.offset += int64(n)
	return n, err
}

func (r *ByteReader) Seek(offset int64, whence int) (abs int64, _ error) {
	if r.Byte == nil {
		return 0, fs.ErrClosed
	}

	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.offset + offset
	case io.SeekEnd:
		abs = int64(r.Len()) + offset
	default:
		return 0, errors.New("buffer.ByteReader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("buffer.ByteReader.Seek: negative position")
	}

	r.offset = abs
	return abs, nil
}

func (r *ByteReader) ReadAt(b []byte, off int64) (n int, err error) {
	if r.Byte == nil {
		return 0, fs.ErrClosed
	}

	if off < 0 {
		return 0, errors.New("buffer.ByteReader.ReadAt: negative offset")
	}
	if off >= int64(r.Len()) {
		return 0, io.EOF
	}
	if n = copy(b, r.B[off:]); n < len(b) {
		err = io.EOF
	}
	return
}

var (
	_ io.Reader   = (*ByteReader)(nil)
	_ io.Seeker   = (*ByteReader)(nil)
	_ io.ReaderAt = (*ByteReader)(nil)
	_ io.Closer   = (*ByteReader)(nil)
	_ Reader      = (*ByteReader)(nil)
)
