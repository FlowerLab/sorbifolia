package parser

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httperr"
)

type testParseResult struct {
	w      [][]byte
	result []byte
	err    error
}

func TestRequestParser_parseMethod(t *testing.T) {
	tests := []testParseResult{
		{[][]byte{[]byte("get / HTTP/1.1")}, []byte("get"), nil},
		{[][]byte{[]byte("GET / HTTP/1.1")}, []byte("GET"), nil},
		{[][]byte{[]byte("HEAD / HTTP/1.1")}, []byte("HEAD"), nil},
		{[][]byte{[]byte("POST / HTTP/1.1")}, []byte("POST"), nil},
		{[][]byte{[]byte("PUT / HTTP/1.1")}, []byte("PUT"), nil},
		{[][]byte{[]byte("PATCH / HTTP/1.1")}, []byte("PATCH"), nil},
		{[][]byte{[]byte("DELETE / HTTP/1.1")}, []byte("DELETE"), nil},
		{[][]byte{[]byte("CONNECT / HTTP/1.1")}, []byte("CONNECT"), nil},
		{[][]byte{[]byte("OPTIONS / HTTP/1.1")}, []byte("OPTIONS"), nil},
		{[][]byte{[]byte("TRACE / HTTP/1.1")}, []byte("TRACE"), nil},
		{[][]byte{[]byte("oneMethod / HTTP/1.1")}, nil, httperr.ParseHTTPMethodErr},

		{[][]byte{[]byte("G"), []byte("ET / HTTP/1.1")}, []byte("GET"), nil},
		{[][]byte{[]byte("GE"), []byte("T / HTTP/1.1")}, []byte("GET"), nil},
		{[][]byte{[]byte("GET"), []byte(" / HTTP/1.1")}, []byte("GET"), nil},
		{[][]byte{[]byte("GET "), []byte("/ HTTP/1.1")}, []byte("GET"), nil},
		{[][]byte{[]byte("GET /"), []byte(" HTTP/1.1")}, []byte("GET"), nil},
		{[][]byte{[]byte("GET / "), []byte("HTTP/1.1")}, []byte("GET"), nil},

		{[][]byte{[]byte("G"), []byte("ET /\r\n")}, []byte("GET"), nil},
		{[][]byte{[]byte("GE"), []byte("T /\r\n")}, []byte("GET"), nil},
		{[][]byte{[]byte("GET"), []byte(" /\r\n")}, []byte("GET"), nil},
		{[][]byte{[]byte("GET "), []byte("/\r\n")}, []byte("GET"), nil},
		{[][]byte{[]byte("GET /"), []byte("\r\n")}, []byte("GET"), nil},

		{[][]byte{[]byte(" / HTTP/1.1")}, nil, nil},
		{[][]byte{[]byte(" /\r\n")}, nil, nil},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("ParseMethod %d", i), func(t *testing.T) {
			var (
				hasCall bool
				err     error
			)
			rp := &RequestParser{SetMethod: func(b []byte) error {
				hasCall = true
				if !bytes.Equal(b, v.result) {
					t.Errorf("in: %v, expected: %v, actual: %v\n", v.w, v.result, b)
				}
				return nil
			}}

			for _, b := range v.w {
				_, err = rp.parseMethod(b)
				if hasCall || err != nil {
					break
				}
			}
			if err != nil {
				if !errors.Is(err, v.err) {
					t.Errorf("in: %v, Err: expected: %v, actual: %v\n", v.w, v.err, err)
				}
			} else if !hasCall {
				t.Errorf("in: %v, expected: %v, actual: none\n", v.w, v.result)
			}

		})
	}
}

func TestRequestParser_parseURI(t *testing.T) {
	tests := []testParseResult{
		{[][]byte{[]byte("/\r\n")}, []byte("/"), nil},
		{[][]byte{[]byte("/ab\r\n")}, []byte("/ab"), nil},
		{[][]byte{[]byte("/23456789-12\r\n")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/1-3-5-7-9-1-3\r\n")}, nil, httperr.RequestURITooLong},
		{[][]byte{[]byte("/"), []byte("23456789-12\r\n")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/23456789"), []byte("-12\r\n")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/23456789-12"), []byte("\r\n")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/23456789-1\r"), []byte("\n")}, []byte("/23456789-1"), nil},

		{[][]byte{[]byte("/23456789-12 HTTP/1.0")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/23456789-12 HTTP/1.1")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/"), []byte("23456789-12 HTTP/1.1")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/23456789"), []byte("-12 HTTP/1.1")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/23456789-12"), []byte(" HTTP/1.1")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/23456789-12 "), []byte("HTTP/1.1")}, []byte("/23456789-12"), nil},
		{[][]byte{[]byte("/23456789-12 HTTP"), []byte("/1.1")}, []byte("/23456789-12"), nil},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("ParseURI %d", i), func(t *testing.T) {
			var (
				hasCall bool
				err     error
			)
			rp := &RequestParser{SetURI: func(b []byte) error {
				hasCall = true
				if !bytes.Equal(b, v.result) {
					t.Errorf("in: %v, expected: %v, actual: %v\n", v.w, v.result, b)
				}
				return nil
			},
				SetVersion: func([]byte) error { return nil },
				Limit:      httpconfig.Config{MaxRequestURISize: 12},
			}

			for _, b := range v.w {
				_, err = rp.parseURI(b)
				if hasCall || err != nil {
					break
				}
			}
			if err != nil {
				if !errors.Is(err, v.err) {
					t.Errorf("in: %v, Err: expected: %v, actual: %v\n", v.w, v.err, err)
				}
			} else if !hasCall {
				t.Errorf("in: %v, expected: %v, actual: none\n", v.w, v.result)
			}

		})
	}
}

func TestRequestParser_parseVersion(t *testing.T) {
	tests := []testParseResult{
		{[][]byte{[]byte("HTTP/1.0\r\n")}, []byte("HTTP/1.0"), nil},
		{[][]byte{[]byte("HTTP/1.1\r\n")}, []byte("HTTP/1.1"), nil},
		{[][]byte{[]byte("HTTP"), []byte("/"), []byte("1.1\r\n")}, []byte("HTTP/1.1"), nil},
		{[][]byte{[]byte("HTTP/"), []byte("1.1\r\n")}, []byte("HTTP/1.1"), nil},
		{[][]byte{[]byte("HTTP/1"), []byte(".1\r\n")}, []byte("HTTP/1.1"), nil},
		{[][]byte{[]byte("HTTP/1.1"), []byte("\r\n")}, []byte("HTTP/1.1"), nil},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var (
				hasCall bool
				err     error
			)
			rp := &RequestParser{SetVersion: func(b []byte) error {
				hasCall = true
				if !bytes.Equal(b, v.result) {
					t.Errorf("in: %v, expected: %v, actual: %v\n", v.w, v.result, b)
				}
				return nil
			},
				Limit: httpconfig.Config{MaxRequestURISize: 12},
			}

			for _, b := range v.w {
				_, err = rp.parseVersion(b)
				if hasCall || err != nil {
					break
				}
			}
			if err != nil {
				if !errors.Is(err, v.err) {
					t.Errorf("in: %v, Err: expected: %v, actual: %v\n", v.w, v.err, err)
				}
			} else if !hasCall {
				t.Errorf("in: %v, expected: %v, actual: none\n", v.w, v.result)
			}

		})
	}
}

func TestRequestParser_parseHeader(t *testing.T) {
	tests := []testParseResult{
		{[][]byte{[]byte("A:b"), []byte("\r\n\r\n")}, []byte("A:b"), nil},
		{[][]byte{[]byte("A:b\r"), []byte("\n\r\n")}, []byte("A:b"), nil},
		{[][]byte{[]byte("A:b\r\n"), []byte("\r\n")}, []byte("A:b"), nil},
		{[][]byte{[]byte("A:b\r\n\r"), []byte("\n")}, []byte("A:b"), nil},
		{[][]byte{[]byte("A:b\r\n\r\n"), []byte("\n")}, []byte("A:b"), nil},

		{[][]byte{[]byte("\r"), []byte("\n\r\n")}, nil, nil},
		{[][]byte{[]byte("\r\n"), []byte("\r\n")}, nil, nil},
		{[][]byte{[]byte("\r\n\r"), []byte("\n")}, nil, nil},
		{[][]byte{[]byte("\r\n\r\n"), []byte("\n")}, nil, nil},

		{[][]byte{[]byte("123456:654321\r\n\r\n")}, nil, httperr.RequestHeaderFieldsTooLarge},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("ParseHeaders %d", i), func(t *testing.T) {
			var (
				hasCall bool
				err     error
			)
			rp := &RequestParser{SetHeaders: func(b []byte) error {
				hasCall = true
				if !bytes.Equal(b, v.result) {
					t.Errorf("in: %v, expected: %v, actual: %v\n", v.w, v.result, b)
				}
				return nil
			}, Limit: httpconfig.Config{MaxRequestHeaderSize: 12}}

			for _, b := range v.w {
				_, err = rp.parseHeader(b)
				if hasCall || err != nil {
					break
				}
			}
			if err != nil {
				if !errors.Is(err, v.err) {
					t.Errorf("in: %v, Err: expected: %v, actual: %v\n", v.w, v.err, err)
				}
			} else if !hasCall {
				t.Errorf("in: %v, expected: %v, actual: none\n", v.w, v.result)
			}
		})
	}
}
