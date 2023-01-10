package httpmessage

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/kv"
	"go.x2ox.com/sorbifolia/http/status"
)

func TestResponse_Read(t *testing.T) {
	tests := []struct {
		res      *Response
		expected []byte
	}{
		{
			res: &Response{
				StatusCode: status.OK,
				Header: httpheader.ResponseHeader{Header: httpheader.Header{KVs: kv.KVs{
					{[]byte("Content-Length"), []byte("12"), false},
				}}},
				Body: strings.NewReader("abc45qwe9012"),
			},
			expected: []byte("200 OK\r\n" +
				"Content-Length: 12\r\n" +
				"\r\n" +
				"abc45qwe9012"),
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := io.ReadAll(v.res)
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(actual, v.expected) {
				t.Errorf("expected %v, got %v", v.expected, actual)
			}
		})
	}
}
