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
