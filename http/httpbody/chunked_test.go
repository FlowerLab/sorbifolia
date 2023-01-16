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
	tests := []*Chunked{
		{
			Data:   make(chan []byte),
			Header: make(chan []byte),
			m:      ModeReadWrite,
			finish: false,
			state:  0,
			once:   sync.Once{},
			buf:    bufpool.Buffer{},
		},
		{
			Data:   make(chan []byte),
			Header: make(chan []byte),
			m:      ModeWrite,
			finish: false,
			state:  chunkedEND,
			once:   sync.Once{},
			buf:    bufpool.Buffer{},
		},
		{
			Data:   make(chan []byte),
			Header: make(chan []byte),
			m:      ModeWrite,
			finish: false,
			once:   sync.Once{},
			buf:    bufpool.Buffer{},
		},
		{
			Data:   make(chan []byte),
			Header: make(chan []byte),
			m:      ModeWrite,
			finish: false,
			once:   sync.Once{},
			buf:    bufpool.Buffer{B: []byte("\r")},
		},
	}

	data := []byte("7\r\nhello, \r\n" +
		"6\r\nworld!\r\n" +
		"0\r\n" +
		"Expires: Fri, 20 Jan 2023 07:28:00 GMT\r\n" +
		"\r\n")
	var wg sync.WaitGroup

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			wg.Add(1)
			go func() {
				if i == 3 {
					data = []byte("\nhello, ")
				}

				_, err := io.Copy(v, bytes.NewReader(data))
				if err != io.EOF && i != 3 {
					t.Error(err)
				}

				if i == 0 || i == 1 || i == 3 {
					close(v.Data)
					close(v.Header)
				}
				wg.Done()
			}()

			var (
				ok  = true
				buf = new(bytes.Buffer)
			)
			for ok {
				_, ok = <-v.Data
			}
			for hv := range v.Header {
				buf.Write(hv)
			}
			t.Log(buf.String())

			v.release()
			wg.Wait()
		})
	}
}
