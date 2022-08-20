package httputils

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"mime/multipart"
	"testing"
)

type EncData struct {
	Foo   string  `json:"foo" xml:"foo"`
	Value float64 `json:"value" xml:"value,omitempty"`
}

var (
	encXML = func() GetEncoder {
		return func(buf io.Writer) Encoder {
			return xml.NewEncoder(buf)
		}
	}

	_encData = EncData{Value: 123.123, Foo: "bar"}
	_enc     = []struct {
		GE     GetEncoder
		Output string
	}{
		{JSON(), "{\"foo\":\"bar\",\"value\":123.123}\n"},
		{encXML(), "<EncData><foo>bar</foo><value>123.123</value></EncData>"},
		{FormData(func(w *multipart.Writer) error {
			_ = w.SetBoundary("go.x2ox.com/sorbifolia/httputils/encoder.FormDataEncoder")
			if err := w.WriteField("foo", "bar"); err != nil {
				return err
			}
			return w.WriteField("value", "123.123")
		}), "--go.x2ox.com/sorbifolia/httputils/encoder.FormDataEncoder\r\n" +
			"Content-Disposition: form-data; name=\"foo\"\r\n\r\n" +
			"bar\r\n--go.x2ox.com/sorbifolia/httputils/encoder.FormDataEncoder\r\n" +
			"Content-Disposition: form-data; name=\"value\"\r\n\r\n123.123\r\n" +
			"--go.x2ox.com/sorbifolia/httputils/encoder.FormDataEncoder--\r\n"},
	}
)

func TestNewEncoder(t *testing.T) {
	var buf = new(bytes.Buffer)

	for _, v := range _enc {
		buf.Reset()

		if err := v.GE(buf).Encode(_encData); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if buf.String() != v.Output {
			t.Errorf("unexpected output: %s", buf.String())
		}
	}
}

func TestFormDataEncoder(t *testing.T) {
	e := FormData(func(w *multipart.Writer) error {
		return errors.New("error")
	})
	var buf = new(bytes.Buffer)

	if err := e(buf); err != nil {
		t.Error("error")
	}
}
