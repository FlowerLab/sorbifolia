package httpbody

import (
	"io"
	"os"

	"go.x2ox.com/sorbifolia/pyrokinesis"
)

var (
	_ io.ReadCloser  = (*TempFile)(nil)
	_ io.WriteCloser = (*TempFile)(nil)
	_ HTTPBody       = (*TempFile)(nil)
)

type TempFile struct {
	*os.File
	filename []byte
	err      error
	mode     rwcMode
}

func (t *TempFile) Filename() string { return pyrokinesis.Bytes.ToString(t.filename) }

func (t *TempFile) Read(p []byte) (n int, err error) {
	if t.File == nil {
		if t.File, err = os.Open(t.Filename()); err != nil {
			return
		}
	}
	return t.File.Read(p)
}

func (t *TempFile) Write(p []byte) (n int, err error) {
	if t.File == nil {
		if t.File, err = os.Open(t.Filename()); err != nil {
			return
		}
	}
	return t.File.Write(p)
}

func (t *TempFile) BodyReader() io.ReadCloser   { return t.getIO(ModeRead) }
func (t *TempFile) BodyWriter() io.WriteCloser  { return t.getIO(ModeWrite) }
func (t *TempFile) getIO(rwc rwcMode) *TempFile { t.mode.SetMode(rwc); return t }

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

func (t *TempFile) create() (err error) {
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
