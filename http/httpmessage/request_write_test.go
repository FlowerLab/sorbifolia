package httpmessage

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"go.x2ox.com/sorbifolia/http/httpbody"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/kv"
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
				ver:    version.Version{Major: 1, Minor: 1},
				Method: "GET",
				Header: httpheader.RequestHeader{
					KVs: kv.KVs{
						{[]byte("Host"), []byte("localhost"), false},
						{[]byte("User-Agent"), []byte("Mozilla/5.0"), false},
						{[]byte("Accept"), []byte("text/html,*/*;q=0.8"), false},
						{[]byte("Accept-Language"), []byte("en-US,en;q=0.3"), false},
						{[]byte("Connection"), []byte("keep-alive"), false},
					},
				},
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
