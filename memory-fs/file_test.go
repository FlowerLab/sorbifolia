package mfs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"reflect"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	tests := []struct {
		*openFile
		offset int64
		whence int
		expect string
		Err    error
	}{
		{
			&openFile{
				file: &file{
					name:    "a.txt",
					modTime: time.Now(),
					data:    []byte("hello,world!"),
				},
				n: 0,
			},

			6,
			io.SeekStart,
			"world!",
			nil,
		},
		{
			&openFile{
				file: &file{
					name:    "a.txt",
					modTime: time.Now(),
					data:    []byte("hello,world!"),
				},
				n: 0,
			},

			-1,
			io.SeekStart,
			"world!",
			&fs.PathError{Op: "seek", Path: "a.txt", Err: fs.ErrInvalid},
		},
		{
			&openFile{
				file: &file{
					name:    "a.txt",
					modTime: time.Now(),
					data:    []byte("hello,world!"),
				},
				n: 0,
			},
			0,
			io.SeekCurrent,
			"hello,world!",
			nil,
		},
		{
			&openFile{
				file: &file{
					name:    "a.txt",
					modTime: time.Now(),
					data:    []byte("hello,world!"),
				},
				n: 0,
			},
			2,
			io.SeekEnd,
			"",
			&fs.PathError{Op: "seek", Path: "a.txt", Err: fs.ErrInvalid},
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, err := v.Seek(v.offset, v.whence)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
			if err == nil {
				var buf bytes.Buffer
				reader := bufio.NewReader(v)
				if _, err = buf.ReadFrom(reader); err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(buf.String(), v.expect) {
					t.Errorf("expect: %v,get: %v", buf.String(), v.expect)
				}
			}
			_ = v.Close()
		})
	}
}
