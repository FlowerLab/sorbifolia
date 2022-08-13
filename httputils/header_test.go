package httputils

import (
	"testing"
)

func TestHTTP_AddHeader(t *testing.T) {
	h := Head().AddHeader("foo", "bar")

	req, _, _ := h.test()
	if string(req.Header.Peek("foo")) != "bar" {
		t.Error("AddHeader err")
	}

	h.AddHeader("foo", "bbb")

	req, _, _ = h.test()
	if string(req.Header.Peek("foo")) != "bar" {
		t.Error("AddHeader err")
	}
}

func TestHTTP_SetHeader(t *testing.T) {
	h := Head().SetHeader("foo", "bar")

	req, _, _ := h.test()
	if string(req.Header.Peek("foo")) != "bar" {
		t.Error("SetHeader err")
	}

	h.SetHeader("foo", "bbb")

	req, _, _ = h.test()
	if string(req.Header.Peek("foo")) != "bbb" {
		t.Error("SetHeader err")
	}
}

func TestHTTP_DelHeader(t *testing.T) {
	h := Head().SetHeader("foo", "bar").DelHeader("foo")

	req, _, _ := h.test()
	if string(req.Header.Peek("foo")) != "" {
		t.Error("DelHeader err")
	}
}

func TestHTTP_SetReferer(t *testing.T) {
	h := Head().SetReferer("foo")

	req, _, _ := h.test()
	if string(req.Header.Referer()) != "foo" {
		t.Error("SetReferer err")
	}
}

func TestHTTP_SetUserAgent(t *testing.T) {
	h := Head().SetUserAgent("foo")

	req, _, _ := h.test()
	if string(req.Header.UserAgent()) != "foo" {
		t.Error("SetUserAgent err")
	}
}

func TestHTTP_SetProtocol(t *testing.T) {
	h := Patch().SetProtocol("foo")

	req, _, _ := h.test()
	if string(req.Header.Protocol()) != "foo" {
		t.Error("SetProtocol err")
	}
}

func TestHTTP_SetByteRange(t *testing.T) {
	h := Head().SetByteRange(1, 123)

	req, _, _ := h.test()
	if string(req.Header.Peek("Range")) != "bytes=1-123" {
		t.Error("SetByteRange err")
	}
}

func TestHTTP_SetContentLength(t *testing.T) {
	h := Put().SetContentLength(123)

	req, _, _ := h.test()
	if req.Header.ContentLength() != 123 {
		t.Error("SetContentLength err")
	}
}

func TestHTTP_SetContentEncoding(t *testing.T) {
	h := Head().SetContentEncoding("br")

	req, _, _ := h.test()
	if string(req.Header.ContentEncoding()) != "br" {
		t.Error("SetContentEncoding err")
	}
}

func TestHTTP_SetMultipartFormBoundary(t *testing.T) {
	h := Head().SetMultipartFormBoundary("go.x2ox.com/sorbifolia/httputils/encoder.FormDataEncoder")

	req, _, _ := h.test()
	if string(req.Header.ContentType()) != "multipart/form-data; boundary=go.x2ox.com/sorbifolia/httputils/encoder.FormDataEncoder" {
		t.Error("SetMultipartFormBoundary err")
	}
}
