package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httperr"
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
	rp := AcquireRequestParser(
		func(b []byte) error { t.actual.method = append(t.actual.method, b...); return nil },
		func(b []byte) error { t.actual.uri = append(t.actual.uri, b...); return nil },
		func(b []byte) error { t.actual.version = append(t.actual.version, b...); return nil },
		func(b []byte) (chunked ChunkedTransfer, err error) {
			t.actual.headers = append(t.actual.headers, b...)
			arr := bytes.Split(t.actual.headers, []byte("\r\n"))

			var setTrailerHeader, setChunked func(b []byte) error = nil, nil
			for _, v := range arr {
				i := bytes.IndexByte(v, ':')
				if i == -1 {
					continue
				}
				if bytes.EqualFold(b[:i], []byte("Transfer-Encoding")) && bytes.Contains(b[i:], []byte("chunked")) {
					setChunked = func(_ []byte) error { return nil }
				}
				if bytes.EqualFold(b[:i], []byte("Trailer")) {
					setTrailerHeader = func(b []byte) error {
						t.actual.chunkedHeaders = append(t.actual.chunkedHeaders, append([]byte{}, b...))
						return nil
					}
				}
			}
			if setChunked != nil {
				chunked = func() (a, b func(b []byte) error) {
					return setTrailerHeader, setChunked
				}
			}

			return
		},
	)

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
					"Connection: keep-alive\r\n" +
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

			fmt.Println(string(v.actual.uri))

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
			rp := &RequestParser{SetHeaders: func(b []byte) (chunked ChunkedTransfer, err error) {
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

func TestRequestParser_parseBodyChunked(t *testing.T) {
	tests := []testParseResult{
		{[][]byte{[]byte("7\r\nMozilla\r\n0\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7"), []byte("\r\nMozilla\r\n0\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r"), []byte("\nMozilla\r\n0\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\n"), []byte("Mozilla\r\n0\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\nM"), []byte("ozilla\r\n0\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\nMozilla"), []byte("\r\n0\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\nMozilla\r"), []byte("\n0\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\nMozilla\r\n"), []byte("0\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\nMozilla\r\n0"), []byte("\r\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\nMozilla\r\n0\r"), []byte("\n\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\nMozilla\r\n0\r\n"), []byte("\r\n")}, []byte("Mozilla"), nil},
		{[][]byte{[]byte("7\r\nMozilla\r\n0\r\r\n"), []byte("\n")}, []byte("Mozilla"), nil},

		{[][]byte{[]byte("0\r\nA:B\r\n\r\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0\r\nA:B\r\n\r"), []byte("\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0\r\nA:B\r\n"), []byte("\r\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0\r\nA:B\r"), []byte("\n\r\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0\r\nA:B"), []byte("\r\n\r\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0\r\nA:"), []byte("B\r\n\r\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0\r\nA"), []byte(":B\r\n\r\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0\r\n"), []byte("A:B\r\n\r\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0\r"), []byte("\nA:B\r\n\r\n")}, []byte("A:B"), nil},
		{[][]byte{[]byte("0"), []byte("\r\nA:B\r\n\r\n")}, []byte("A:B"), nil},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var hasCall bool

			rp := &RequestParser{
				state: ReadBodyChunked,
				setTrailerHeader: func(b []byte) error {
					hasCall = true
					if !bytes.Equal(b, v.result) {
						t.Errorf("in: %v, expected: %v, actual: %v\n", v.w, v.result, b)
					}
					return nil
				},
				setChunked: func(b []byte) error {
					hasCall = true
					if !bytes.Equal(b, v.result) {
						t.Errorf("in: %v, expected: %v, actual: %v\n", v.w, v.result, b)
					}
					return nil
				},
			}

			for _, b := range v.w {
				var length = len(b)
				for rn := 0; rn < length; {
					n, err := rp.parseBodyChunked(b[rn:])
					if err != nil {
						if !errors.Is(err, v.err) {
							t.Errorf("in: %v, Err: expected: %v, actual: %v\n", v.w, v.err, err)
						}
						return
					}
					rn += n
				}
			}

			if !hasCall {
				t.Errorf("in: %v, expected: %v, actual: none\n", v.w, v.result)
			}
		})
	}
}
