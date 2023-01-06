package http

import (
	"bytes"
	"fmt"
	"io"

	"go.x2ox.com/sorbifolia/http/httpbody"
	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/parser"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

type Request struct {
	ver    version.Version
	Method method.Method
	Header RequestHeader
	Body   io.ReadCloser
}

func (r *Request) parse(read io.Reader) {
	p := parser.AcquireRequestParser(
		func(b []byte) error { r.Method = method.Parse(b); return nil },
		func(b []byte) error { r.Header.RequestURI = append(r.Header.RequestURI, b...); return nil },
		func(b []byte) error {
			var ok bool
			if r.ver, ok = version.Parse(b); !ok {
				return httperr.ParseHTTPVersionErr
			}
			return nil
		},
		func(b []byte) (chunked parser.ChunkedTransfer, length int, err error) {
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
			return
		},
	)
	m := &httpbody.Memory{}
	p.BW = m.BodyWriter()
	r.Body = m.BodyReader()
	if _, err := util.Copy(p, read); err != nil && err != io.EOF {
		fmt.Println(err)
	}
}
