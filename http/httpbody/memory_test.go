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

	memo.release()
}

type MER struct {
	*Memory
	Err error
}

func TestMemoryErrorRead(t *testing.T) {
	test := []MER{
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

	for i, v := range test {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			buf := new(bytes.Buffer)
			_, err := io.Copy(buf, v.Memory)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expected: %v,got: %v", v.Err, err)
			}
		})
	}
}
