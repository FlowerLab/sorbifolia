package httpbody

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
)

func TestMemoryRead(t *testing.T) {
	data := []byte("hello,world")
	memo := AcquireMemory()
	memo.mode = ModeReadWrite
	memo.buf.B = data

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, memo); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(string(data), buf.String()) {
		t.Errorf("expected: %v got: %v", string(data), buf.String())
	}

	_ = memo.Close()
	memo.release()
}

func TestMemoryErrorRead(t *testing.T) {
	tests := []MER{
		{
			Memory: &Memory{
				buf:  bufpool.Buffer{},
				p:    0,
				mode: ModeRead,
			},
			Err: nil,
		},
		{
			Memory: &Memory{
				buf:  bufpool.Buffer{B: []byte("hello,world")},
				p:    0,
				mode: ModeWrite,
			},
			Err: httperr.ErrNotYetReady,
		},
		{
			Memory: &Memory{
				buf:  bufpool.Buffer{B: []byte("hello,world")},
				p:    0,
				mode: ModeClose,
			},
			Err: nil,
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			_, err := io.Copy(buf, v.Memory)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expected: %v,got: %v", v.Err, err)
			}
			_ = v.Memory.Close()
		})
	}
}

func TestMemoryWrite(t *testing.T) {
	tests := []MW{
		{
			Memory: &Memory{
				buf:  bufpool.Buffer{},
				p:    0,
				mode: ModeReadWrite,
			},
			Data: []byte("11\r\nhello,world\r\n"),
			Res:  "11\r\nhello,world\r\n",
			Err:  nil,
		},
		{
			Memory: &Memory{
				buf:  bufpool.Buffer{},
				p:    0,
				mode: ModeWrite,
			},
			Data: []byte("123456"),
			Res:  "123456",
			Err:  nil,
		},
		{
			Memory: &Memory{
				buf:  bufpool.Buffer{},
				p:    0,
				mode: ModeRead,
			},
			Data: []byte("123456"),
			Res:  "",
			Err:  io.EOF,
		},
		{
			Memory: &Memory{
				buf:  bufpool.Buffer{},
				p:    0,
				mode: ModeClose,
			},
			Data: []byte("123456"),
			Res:  "",
			Err:  io.EOF,
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, err := io.Copy(v.Memory, bytes.NewReader(v.Data))
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expected: %v,got: %v", v.Err, err)
			}
			if !reflect.DeepEqual(v.Res, v.buf.String()) {
				t.Errorf("expected: %v,got: %v", v.Res, v.buf.String())
			}
		})
	}
}

type MER struct {
	*Memory
	Err error
}

type MW struct {
	*Memory
	Data []byte
	Res  string
	Err  error
}
