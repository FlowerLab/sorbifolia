package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/httpheader"
)

type testRequestParserWriteResult struct {
	method, uri, version, headers, body []byte
	chunkedHeaders                      [][]byte
}

type testRequestParserWrite struct {
	r        io.Reader
	expected testRequestParserWriteResult
	actual   testRequestParserWriteResult
}

func (t *testRequestParserWrite) genRequestParser() *RequestParser {
	rp := AcquireRequestParser()
	rp.SetMethod = func(b []byte) error { t.actual.method = append(t.actual.method, b...); return nil }
	rp.SetURI = func(b []byte) error { t.actual.uri = append(t.actual.uri, b...); return nil }
	rp.SetVersion = func(b []byte) error { t.actual.version = append(t.actual.version, b...); return nil }
	rp.SetHeaders = func(b []byte) (length int, err error) {
		t.actual.headers = append(t.actual.headers, b...)
		arr := bytes.Split(t.actual.headers, []byte("\r\n"))

		var isChunked bool
		for _, v := range arr {
			i := bytes.IndexByte(v, ':')
			if i == -1 {
				continue
			}
			if bytes.EqualFold(b[:i], []byte("Transfer-Encoding")) && bytes.Contains(b[i:], []byte("chunked")) {
				isChunked = true
			}
			if bytes.EqualFold(b[:i], []byte("Content-Length")) && bytes.Contains(b[i:], []byte("chunked")) {
				isChunked = true
				i++
				for b[i] == ' ' {
					i++
				}
				length = int(httpheader.ContentLength(b[i:]).Length())
			}
		}

		if isChunked {
			length = -1
		}
		return
	}

	return rp
}

func TestRequestParser_Write(t *testing.T) {
	tests := []testRequestParserWrite{
		{
			r: strings.NewReader(
				"GET / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"User-Agent: Mozilla/5.0\r\n" +
					"Accept: text/html,*/*;q=0.8\r\n" +
					"Accept-Language: en-US,en;q=0.3\r\n" +
					"Connection: keep-alive\r\n\r\n" +
					""),
			expected: testRequestParserWriteResult{
				method: []byte("GET"), uri: []byte("/"), version: []byte("HTTP/1.1"),
				headers: []byte("Host: localhost\r\n" +
					"User-Agent: Mozilla/5.0\r\n" +
					"Accept: text/html,*/*;q=0.8\r\n" +
					"Accept-Language: en-US,en;q=0.3\r\n" +
					"Connection: keep-alive"),
				body: nil, chunkedHeaders: nil,
			},
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			rp := v.genRequestParser()
			defer ReleaseRequestParser(rp)

			if _, err := io.Copy(rp, v.r); err != nil && !errors.Is(err, io.EOF) {
				t.Error(err)
			}

			if !reflect.DeepEqual(v.expected, v.actual) {
				t.Errorf("expected: %v, actual: %v\n", v.expected, v.actual)
			}
		})
	}
}

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
			rp := &RequestParser{SetHeaders: func(b []byte) (length int, err error) {
				hasCall = true
				if !bytes.Equal(b, v.result) {
					t.Errorf("in: %v, expected: %v, actual: %v\n", v.w, v.result, b)
				}
				return
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
