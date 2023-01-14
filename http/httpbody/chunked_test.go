package httpbody

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sync"
	"testing"
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
	hex_data := fmt.Sprintf("%x", len(data))

	expect := []byte(hex_data +
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
	chunked.release()
}

func TestWrite(t *testing.T) {
	data := []byte("7\r\n1234567\r\n" +
		"5\r\n25380\r\n" +
		"0\r\n" +
		"Expires: Fri, 20 Jan 2023 07:28:00 GMT\r\n" +
		"\r\n")

	wc := AcquireChunked()
	wc.m = ModeWrite
	wc.Data = make(chan []byte)
	wc.Header = make(chan []byte)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		v, ok := <-wc.Data
		for ok {
			t.Log(string(v))
			v, ok = <-wc.Data
		}

		buf := new(bytes.Buffer)
		for v = range wc.Header {
			buf.Write(v)
		}
		t.Log(buf.String())

		wg.Done()
	}()

	_, err := io.Copy(wc, bytes.NewReader(data))
	if err != io.EOF {
		t.Error(err)
	}

	wg.Wait()
	wc.release()
}
