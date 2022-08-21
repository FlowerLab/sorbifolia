package httputils

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestJSONParser(t *testing.T) {
	var (
		m    = make(map[string]string)
		resp = &fasthttp.Response{}
	)
	resp.SetBody([]byte(`{"a":"a"}`))

	if err := JSONParser(&m)(resp); err != nil {
		t.Error("err")
	}

	resp.Header.SetContentEncoding("xxx")
	if err := JSONParser(m)(resp); err == nil {
		t.Error("err")
	}
}

func TestHeaderParser(t *testing.T) {
	var (
		m    = make(map[string]string)
		resp = &fasthttp.Response{}
	)
	resp.Header.Set("a", "b")
	resp.Header.Set("b", "a")
	resp.Header.Set("c", "c")
	m["a"] = ""
	m["b"] = ""

	if err := HeaderParser(m)(resp); err != nil {
		t.Error("err")
	}
	for _, v := range m {
		if !(v == "a" || v == "b") {
			t.Error("expected")
		}
	}
}
