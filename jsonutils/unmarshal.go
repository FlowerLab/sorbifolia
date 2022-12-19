package jsonutils

import (
	"bytes"
	"encoding/json"
	"sync"
)

type decoder struct {
	dec *json.Decoder
	r   *bytes.Reader
}

type SyntaxError struct {
	msg    string
	Offset int64
}

func (e *SyntaxError) Error() string { return e.msg }

var (
	readerPool  = sync.Pool{New: func() any { return bytes.NewReader(nil) }}
	decoderPool = sync.Pool{New: func() any {
		var d decoder
		d.r = readerPool.Get().(*bytes.Reader)
		d.dec = json.NewDecoder(d.r)
		return &d
	}}
)

func Unmarshal(b []byte, v any) error {
	d := decoderPool.Get().(*decoder)
	d.r.Reset(b)
	off := d.dec.InputOffset()
	err := d.dec.Decode(v)
	d.r.Reset(nil)

	switch je := err.(type) {
	case *json.SyntaxError:
		je.Offset -= off
	case *json.UnmarshalTypeError:
		je.Offset -= off
	case nil:
		if d.dec.More() {
			err = &SyntaxError{
				Offset: d.dec.InputOffset() - off,
				msg:    "trailing garbage; see https://github.com/golang/go/issues/36225",
			}
		}
	}

	if err == nil {
		decoderPool.Put(d)
	} else {
		readerPool.Put(d.r)
	}
	return err
}
