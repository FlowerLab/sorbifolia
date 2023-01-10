package httpmessage

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"go.x2ox.com/sorbifolia/http/httpbody"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/kv"
	"go.x2ox.com/sorbifolia/http/url"
	"go.x2ox.com/sorbifolia/http/version"
)

func TestRequest_Write(t *testing.T) {
	tests := []struct {
		data     []byte
		expected *Request
	}{
		{
			data: []byte("GET / HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"User-Agent: Mozilla/5.0\r\n" +
				"Accept: text/html,*/*;q=0.8\r\n" +
				"Accept-Language: en-US,en;q=0.3\r\n" +
				"Connection: keep-alive\r\n\r\n" +
				""),
			expected: &Request{
				Version: version.Version{Major: 1, Minor: 1},
				Method:  "GET",
				Header: httpheader.RequestHeader{
					Header: httpheader.Header{KVs: kv.KVs{
						{K: []byte("Host"), V: []byte("localhost")},
						{K: []byte("User-Agent"), V: []byte("Mozilla/5.0")},
						{K: []byte("Accept"), V: []byte("text/html,*/*;q=0.8")},
						{K: []byte("Accept-Language"), V: []byte("en-US,en;q=0.3")},
						{K: []byte("Connection"), V: []byte("keep-alive")},
					}}},
				Body: httpbody.Null(),
			},
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r := &Request{}
			if _, err := io.Copy(r, bytes.NewReader(v.data)); err != nil && err != io.EOF {
				t.Error(err)
			}

			if !reflect.DeepEqual(r.Header.KVs, v.expected.Header.KVs) {
				t.Errorf("expected %v, got %v", v.expected.Header.KVs, r.Header.KVs)
			}
			if !reflect.DeepEqual(r.Method, v.expected.Method) {
				t.Errorf("expected %v, got %v", v.expected.Header.KVs, r.Header.KVs)
			}
		})
	}
}

func TestRequest_Read(t *testing.T) {
	req := &Request{
		Version: version.Version{Major: 1, Minor: 1},
		Method:  "",
		Header: httpheader.RequestHeader{
			Header: httpheader.Header{KVs: kv.KVs{
				kv.KV{K: char.ContentLength, V: []byte("18")},
			}},
			RemoteAddr: nil,
			RequestURI: nil,
			URL: url.URL{
				Scheme:   []byte("https"),
				Host:     []byte("example.com"),
				Path:     []byte("/index.html"),
				Query:    []byte("id=1"),
				Fragment: []byte("post"),
			},
		},
		Body: NopWriteCloser{Reader: strings.NewReader("111-222-333-444-18")},
	}
	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, req); err != nil {
		t.Error(err)
	}
	fmt.Println(buf.String())
}

type NopWriteCloser struct {
	io.Reader
}

func (n2 NopWriteCloser) Write(p []byte) (n int, err error) { return 0, io.EOF }
func (n2 NopWriteCloser) Close() error                      { return nil }
