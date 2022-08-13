package httputils

import (
	"testing"
)

func TestNewContentType(t *testing.T) {
	if NewContentType("foo") != "foo" {
		t.Errorf("NewContentType err")
	}
}

func TestContentType_SetCharset(t *testing.T) {
	if TextHTML.SetCharset("utf-8") != "text/html; charset=utf-8" {
		t.Error("SetCharset err")
	}
	if TextPlain.SetCharset("UTF-8") != "text/plain; charset=UTF-8" {
		t.Error("SetCharset err")
	}
	if AppOctetStream.SetCharset("UTF-16") != "application/octet-stream; charset=UTF-16" {
		t.Error("SetCharset err")
	}
	if AppFormUrlencoded.SetCharset("UTF-32") !=
		"application/x-www-form-urlencoded; charset=UTF-32" {
		t.Error("SetCharset err")
	}
	if AppXML.SetCharset("GBK") != "application/xml; charset=GBK" {
		t.Error("SetCharset err")
	}
}

func TestContentType_SetBoundary(t *testing.T) {
	if MultiFormData.SetBoundary("go.x2ox.com/sorbifolia/httputils/encoder.FormDataEncoder") !=
		"multipart/form-data; boundary=go.x2ox.com/sorbifolia/httputils/encoder.FormDataEncoder" {
		t.Error("SetBoundary err")
	}
}

func TestHTTP_SetContentType(t *testing.T) {
	h := Trace().SetContentType(AppJSON)
	req, _, _ := h.test()
	if string(req.Header.ContentType()) != "application/json" {
		t.Error("SetContentType err")
	}
}
