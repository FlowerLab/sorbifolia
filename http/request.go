package http

import (
	"bytes"
	"fmt"
	"io"

	"go.x2ox.com/sorbifolia/http/httpbody"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/parser"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

type Request struct {
	ver    version.Version
	Method method.Method
	Header httpheader.RequestHeader
	Body   io.ReadCloser
}

func (r *Request) parse(read io.Reader) {
	p := parser.AcquireRequestWriter()
	p.SetMethod = func(b []byte) error { r.Method = method.Parse(b); return nil }
	p.SetURI = func(b []byte) error { r.Header.RequestURI = append(r.Header.RequestURI, b...); return nil }
	p.SetVersion = func(b []byte) error {
		var ok bool
		if r.ver, ok = version.Parse(b); !ok {
			return httperr.ParseHTTPVersionErr
		}
		return nil
	}
	p.SetHeaders = func(b []byte) (length int, err error) {
		r.Header.KVs.PreAlloc(bytes.Count(b, char.CRLF) + 1)

		for len(b) > 0 {
			idx := bytes.Index(b, char.CRLF)
			switch idx {
			case 0:
			case -1:
				r.Header.KVs.AddHeader(b)
				idx = len(b) - 2
			default:
				r.Header.KVs.AddHeader(b[:idx])
			}

			b = b[idx+2:]
		}

		if err = r.Header.RawParse(); err != nil {
			return
		}
		length = int(r.Header.ContentLength.Length())

		r.Header.TransferEncoding.Each(func(val []byte) bool {
			if bytes.Equal(val, char.Chunked) {
				length = -1
				return false
			}
			return true
		})

		switch length {
		case 0:
			r.Body = httpbody.Null()
		case -1:
			c := &httpbody.Chunked{
				Data:   make(chan []byte, 1),
				Header: make(chan []byte, 1),
			}
			p.BW = c
			r.Body = c
		default:
			m := &httpbody.Memory{}
			p.BW = m
			r.Body = m
		}
		if length == 0 {
			r.Body = httpbody.Null()
		}

		return
	}

	if _, err := util.Copy(p, read); err != nil && err != io.EOF {
		fmt.Println(err)
	}
}
