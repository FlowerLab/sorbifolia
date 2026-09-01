package jsonutils

import (
	"bytes"
	"encoding/json"
	"sync"
)

type SyntaxError struct {
	msg    string
	Offset int64
}

func (e *SyntaxError) Error() string { return e.msg }

var readerPool = sync.Pool{New: func() any { return bytes.NewReader(nil) }}

func Unmarshal(b []byte, v any) error {
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(b)
	dec := json.NewDecoder(r)
	err := dec.Decode(v)
	if err == nil && dec.More() {
		err = &SyntaxError{
			Offset: dec.InputOffset(),
			msg:    "trailing garbage; see https://github.com/golang/go/issues/36225",
		}
	}
	r.Reset(nil)
	readerPool.Put(r)
	return err
}
