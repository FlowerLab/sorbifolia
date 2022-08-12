package httputils

import (
	"encoding/json"
	"io"
	"mime/multipart"
)

type Encoder interface {
	Encode(v any) error
}

type GetEncoder func(buf io.Writer) Encoder

func JSON() GetEncoder {
	return func(buf io.Writer) Encoder {
		return json.NewEncoder(buf)
	}
}

func FormData(fn func(*multipart.Writer) error) GetEncoder {
	return func(buf io.Writer) Encoder {
		return FormDataEncoder{w: buf, fn: fn}
	}
}

type FormDataEncoder struct {
	w  io.Writer
	fn func(w *multipart.Writer) error
}

func (f FormDataEncoder) Encode(_ any) error {
	mw := multipart.NewWriter(f.w)
	if err := f.fn(mw); err != nil {
		return err
	}
	return mw.Close()
}
