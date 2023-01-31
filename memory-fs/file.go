package mfs

import (
	"io"
	"io/fs"
	"time"
)

type file struct {
	name    string
	perm    fs.FileMode
	modTime time.Time
	data    []byte
}

type openFile struct {
	*file
	n int64
}

func (f *file) FS() fs.File                { return &openFile{file: f} }
func (f *file) Name() string               { return f.name }
func (f *file) Size() int64                { return int64(len(f.data)) }
func (f *file) Mode() fs.FileMode          { return f.perm }
func (f *file) ModTime() time.Time         { return f.modTime }
func (f *file) IsDir() bool                { return false }
func (f *file) Sys() any                   { return nil }
func (f *file) Stat() (fs.FileInfo, error) { return f, nil }
func (f *file) Info() (fs.FileInfo, error) { return f, nil }
func (f *file) Type() fs.FileMode          { return f.perm }
func (f *file) Close() error               { return nil }

func (f *openFile) Read(b []byte) (int, error) {
	if f.n == f.Size() {
		return 0, io.EOF
	}

	n := copy(b, f.data[f.n:])
	f.n += int64(n)
	return n, nil
}

func (f *openFile) Seek(offset int64, whence int) (int64, error) {
	size := f.Size()

	switch whence {
	case io.SeekStart:
	case io.SeekCurrent:
		offset += f.n
	case io.SeekEnd:
		offset += size
	}
	if offset < 0 || offset > size {
		return 0, &fs.PathError{Op: "seek", Path: f.Name(), Err: fs.ErrInvalid}
	}
	f.n = offset
	return offset, nil
}

var (
	_ io.Seeker = (*openFile)(nil)
	_ fs.File   = (*openFile)(nil)
	_ openFS    = (*file)(nil)
)
