package httpbody

import (
	"io"
	"os"

	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type TempFile struct {
	*os.File
	filename []byte
	err      error
	mode     rwcMode
}

func (t *TempFile) Read(p []byte) (n int, err error) {
	switch t.mode {
	case ModeReadWrite:
		t.mode = ModeRead
		if t.File, err = os.Open(t.Filename()); err != nil {
			return
		}
	case ModeRead:
	case ModeWrite, ModeClose:
		return 0, io.EOF
	default:
		panic("BUG: unknown state")
	}

	return t.File.Read(p)
}

func (t *TempFile) Write(p []byte) (n int, err error) {
	switch t.mode { // 0:rw, 1:w, 2:r, 3:c
	case ModeReadWrite:
		t.mode = ModeWrite
		if t.File, err = os.OpenFile(t.Filename(), os.O_WRONLY, 0777); err != nil {
			return
		}
	case ModeWrite:
	case ModeRead, ModeClose:
		return 0, io.EOF
	default:
		panic("BUG: unknown state")
	}

	return t.File.Write(p)
}

func (t *TempFile) Close() error {
	switch t.mode { // 0:rw, 1:w, 2:r, 3:c
	case ModeReadWrite:
		t.mode = ModeClose
		if err := os.Remove(t.Filename()); err != nil {
			return err
		}
	case ModeRead, ModeWrite:
		t.mode = ModeReadWrite
		if err := t.File.Close(); err != nil {
			return err
		}
		t.File = nil
	case ModeClose:
	default:
		panic("BUG: unknown state")
	}

	return nil
}

func (t *TempFile) Reset() {
	t.File = nil
	t.filename = t.filename[:0]
	t.err = nil
	t.mode = ModeReadWrite
}

func (t *TempFile) release()         { t.Reset(); _TempFilePool.Put(t) }
func (t *TempFile) Filename() string { return pyrokinesis.Bytes.ToString(t.filename) }

func (t *TempFile) Init() (err error) {
	if t.File, err = os.CreateTemp("", tempFilePattern); err != nil {
		return err
	}
	if err = t.File.Close(); err != nil {
		return err
	}
	t.filename = append(t.filename, t.File.Name()...)
	t.File = nil
	return
}

const tempFilePattern = "sorbifolia-http-tmp-*"
