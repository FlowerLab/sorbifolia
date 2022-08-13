package httputils

import (
	"bufio"
	"strings"
	"testing"
)

func TestHTTP_AppendBody(t *testing.T) {
	h := Get().AppendBody([]byte("hello")).AppendBodyString(" world")
	req, _, _ := h.test()
	if string(req.Body()) != "hello world" {
		t.Error("append body failed")
	}
}

func TestHTTP_AppendBodyString(t *testing.T) {
	h := Get().AppendBody([]byte("hello")).AppendBodyString(" world")
	req, _, _ := h.test()
	if string(req.Body()) != "hello world" {
		t.Error("append body failed")
	}
}

func TestHTTP_ReadBody(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	h := Get().ReadBody(r, 5, 123)
	req, _, _ := h.test()
	if string(req.Body()) != "hello" {
		t.Error("read body failed")
	}
}

func TestHTTP_SetBody(t *testing.T) {
	h := Get().SetBody([]byte("hello"))
	req, _, _ := h.test()
	if string(req.Body()) != "hello" {
		t.Error("set body failed")
	}
}

func TestHTTP_SetBodyStream(t *testing.T) {
	r := strings.NewReader("hello world")

	h := Get().SetBodyStream(r, 11)
	req, _, _ := h.test()
	if string(req.Body()) != "hello world" {
		t.Error("set body stream failed")
	}
}

func TestHTTP_SetBodyString(t *testing.T) {
	h := Get().SetBodyString("hello world")
	req, _, _ := h.test()
	if string(req.Body()) != "hello world" {
		t.Error("set body string failed")
	}
}

func TestHTTP_SetConnectionClose(t *testing.T) {
	h := Get().SetConnectionClose()
	req, _, _ := h.test()
	if !req.Header.ConnectionClose() {
		t.Error("SetConnectionClose failed")
	}
}

func TestHTTP_SetRequestURI(t *testing.T) {
	h := Get().SetRequestURI("https://ip.x2ox.com")
	req, _, _ := h.test()
	if req.URI().String() != "https://ip.x2ox.com/" {
		t.Errorf("SetRequestURI failed %s", req.URI().String())
	}
}

func TestHTTP_SetHost(t *testing.T) {
	h := Get().SetHost("https://ip.x2ox.com")
	req, _, _ := h.test()
	if string(req.Host()) != "https://ip.x2ox.com" {
		t.Errorf("SetHost failed %s", string(req.Host()))
	}
}
