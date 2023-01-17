package httpbody

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sync"
	"testing"

	"go.x2ox.com/sorbifolia/http/internal/bufpool"
)

func TestRead(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	chunked := AcquireChunked()
	chunked.m = ModeRead
	chunked.Data = make(chan []byte)
	chunked.Header = make(chan []byte)

	data := []byte("7\r\nhello, \r\n6\r\nworld!\r\n")
	header := []byte("Expires: Fri, 20 Jan 2023 07:28:00 GMT\r\n")
	hexData := fmt.Sprintf("%x", len(data))

	expect := []byte(hexData +
		"\r\n" +
		"7\r\nhello, \r\n6\r\nworld!\r\n" +
		"0\r\n" +
		"Expires: Fri, 20 Jan 2023 07:28:00 GMT\r\n" +
		"\r\n")
	go func() {
		buf := new(bytes.Buffer)
		_, err := io.Copy(buf, chunked)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(expect, buf.Bytes()) {
			t.Errorf("expectd: %v,got: %v", string(expect), buf.String())
		}

		wg.Done()
	}()

	chunked.Data <- data
	close(chunked.Data)
	chunked.Header <- header
	close(chunked.Header)

	wg.Wait()

	_ = chunked.Close()
	chunked.release()
}

func TestWrite(t *testing.T) {
	tests := []TC{
		{
			Chunked: &Chunked{
				Data:   make(chan []byte),
				Header: make(chan []byte),
				m:      ModeReadWrite,
				finish: false,
				state:  0,
				once:   sync.Once{},
				buf:    bufpool.Buffer{},
			},
			Data: []byte("7\r\nhello, \r\n" +
				"6\r\nworld!\r\n" +
				"0\r\n" +
				"Expires: Fri, 20 Jan 2023 07:28:00 GMT\r\n" +
				"\r\n"),
			Res: "",
			Fn: func(chunked *Chunked) {
				close(chunked.Data)
				close(chunked.Header)
			},
			Err: io.EOF,
		},
		{
			Chunked: &Chunked{
				Data:   make(chan []byte),
				Header: make(chan []byte),
				m:      ModeWrite,
				finish: false,
				state:  chunkedEND,
				once:   sync.Once{},
				buf:    bufpool.Buffer{},
			},
			Data: []byte("7\r\nhello, \r\n" +
				"6\r\nworld!\r\n" +
				"0\r\n" +
				"Expires: Fri, 20 Jan 2023 07:28:00 GMT\r\n" +
				"\r\n"),
			Res: "",
			Fn: func(chunked *Chunked) {
				close(chunked.Data)
				close(chunked.Header)
			},
			Err: io.EOF,
		},
		{
			Chunked: &Chunked{
				Data:   make(chan []byte),
				Header: make(chan []byte),
				m:      ModeWrite,
				finish: false,
				once:   sync.Once{},
				buf:    bufpool.Buffer{},
			},
			Data: []byte("7\r\nhello, \r\n" +
				"6\r\nworld!\r\n" +
				"0\r\n" +
				"Expires: Fri, 20 Jan 2023 07:28:00 GMT\r\n" +
				"\r\n"),
			Res: "Expires: Fri, 20 Jan 2023 07:28:00 GMT",
			Fn:  nil,
			Err: io.EOF,
		},
		{

			Chunked: &Chunked{
				Data:   make(chan []byte),
				Header: make(chan []byte),
				m:      ModeWrite,
				finish: false,
				once:   sync.Once{},
				buf:    bufpool.Buffer{B: []byte("\r")},
			},
			Data: []byte("\nhello, "),
			Res:  "",
			Fn: func(chunked *Chunked) {
				close(chunked.Data)
				close(chunked.Header)
			},
			Err: nil,
		},
	}

	var wg sync.WaitGroup
	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			wg.Add(1)
			go func() {
				_, err := io.Copy(v.Chunked, bytes.NewReader(v.Data))
				if !reflect.DeepEqual(err, v.Err) {
					t.Errorf("expected: %v,got: %v", v.Err, err)
				}
				if v.Fn != nil {
					v.Fn(v.Chunked)
				}
				wg.Done()
			}()

			var (
				ok  = true
				buf = new(bytes.Buffer)
			)
			for ok {
				_, ok = <-v.Chunked.Data
			}
			for hv := range v.Header {
				buf.Write(hv)
			}

			if !reflect.DeepEqual(v.Res, buf.String()) {
				t.Errorf("expected: %v,got: %v", v.Res, buf.String())
			}
			wg.Wait()
		})
	}
}

type TC struct {
	*Chunked
	Data []byte
	Res  string
	Fn   func(c *Chunked)
	Err  error
}
