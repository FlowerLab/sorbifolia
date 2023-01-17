package httpbody

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestTempFileRead(t *testing.T) {
	var (
		err  error
		f    *os.File
		data = []byte("hello,123")
	)

	tf := AcquireTempFile()
	tf.mode = ModeRead
	if err = tf.Init(); err != nil {
		t.Error(err)
	}

	f, err = os.OpenFile(tf.Filename(), os.O_RDWR, 0777)
	if err != nil {
		t.Error(err)
	}
	_, err = f.Write(data)
	_ = f.Close()

	tf.File, err = os.Open(tf.Filename())
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, tf); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(string(data), buf.String()) {
		t.Errorf("expected: %v,got: %v", string(data), buf.String())
	}

	tf.File.Close()
	_ = os.Remove(tf.File.Name())
	tf.release()
}

func TestTempFileErrorRead(t *testing.T) {
	var (
		f        *os.File
		err      error
		filename = "text.txt"
		data     = []byte("w123456asd")
	)

	if f, err = os.Create(filename); err != nil {
		t.Error(err)
	}
	if _, err = f.Write(data); err != nil {
		t.Error(err)
	}
	f.Close()

	tests := []TfER{
		{
			TempFile: &TempFile{
				File:     nil,
				filename: []byte(filename),
				err:      nil,
				mode:     ModeReadWrite,
			},
			Err: nil,
			Res: string(data),
		},
		{
			TempFile: &TempFile{
				File:     nil,
				filename: nil,
				err:      nil,
				mode:     ModeWrite,
			},
			Err: nil,
			Res: "",
		},
		{
			TempFile: &TempFile{
				File:     nil,
				filename: nil,
				err:      nil,
				mode:     ModeClose,
			},
			Err: nil,
			Res: "",
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, v.TempFile)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expected: %v,got: %v", v.Err, err)
			}
			if !reflect.DeepEqual(v.Res, buf.String()) {
				t.Errorf("expected: %v,got: %v", v.Res, buf.String())
			}
			_ = v.TempFile.Close()
		})
	}

	if err = os.Remove(filename); err != nil {
		t.Error(err)
	}
}

func TestTempFileWrite(t *testing.T) {
	var (
		f        *os.File
		err      error
		filename = "text1.txt"
		data     = []byte("w123456asd")
		p        []byte
	)

	if f, err = os.Create(filename); err != nil {
		t.Error(err)
	}
	f.Close()

	tests := []TfER{
		{
			TempFile: &TempFile{
				File:     nil,
				filename: []byte(filename),
				err:      nil,
				mode:     ModeReadWrite,
			},
			Err: nil,
			Res: "",
		},
		{
			TempFile: &TempFile{
				File:     nil,
				filename: []byte(filename),
				err:      nil,
				mode:     ModeRead,
			},
			Err: io.EOF,
			Res: "",
		},
		{
			TempFile: &TempFile{
				File:     nil,
				filename: []byte(filename),
				err:      nil,
				mode:     ModeClose,
			},
			Err: io.EOF,
			Res: "",
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, err = io.Copy(v.TempFile, bytes.NewReader(data))
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expected: %v,got: %v", v.Err, err)
			}
			if v.TempFile.File != nil {
				p, err = os.ReadFile(filename)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(data, p) {
					t.Errorf("expected: %v,got: %v", string(data), string(p))
				}
			}
			_ = v.TempFile.Close()
		})
	}

	if err = os.Remove(filename); err != nil {
		t.Error(err)
	}
}

type TfER struct {
	*TempFile
	Err error
	Res string
}
