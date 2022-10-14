package httputils

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"
	"testing"
)

func TestHTTP_AppendBody(t *testing.T) {
	t.Parallel()

	h := Get().AppendBody([]byte("hello")).AppendBodyString(" world")
	req, _, _ := h.test()
	if string(req.Body()) != "hello world" {
		t.Error("append body failed")
	}
}

func TestHTTP_AppendBodyString(t *testing.T) {
	t.Parallel()

	h := Get().AppendBody([]byte("hello")).AppendBodyString(" world")
	req, _, _ := h.test()
	if string(req.Body()) != "hello world" {
		t.Error("append body failed")
	}
}

func TestHTTP_ReadBody(t *testing.T) {
	t.Parallel()

	r := bufio.NewReader(strings.NewReader("hello world"))
	h := Get().ReadBody(r, 5, 123)
	req, _, _ := h.test()
	if string(req.Body()) != "hello" {
		t.Error("read body failed")
	}
}

func TestHTTP_SetBody(t *testing.T) {
	t.Parallel()

	h := Get().SetBody([]byte("hello"))
	req, _, _ := h.test()
	if string(req.Body()) != "hello" {
		t.Error("set body failed")
	}
}

func TestHTTP_SetBodyStream(t *testing.T) {
	t.Parallel()

	r := strings.NewReader("hello world")

	h := Get().SetBodyStream(r, 11)
	req, _, _ := h.test()
	if string(req.Body()) != "hello world" {
		t.Error("set body stream failed")
	}
}

func TestHTTP_SetBodyString(t *testing.T) {
	t.Parallel()

	h := Get().SetBodyString("hello world")
	req, _, _ := h.test()
	if string(req.Body()) != "hello world" {
		t.Error("set body string failed")
	}
}

func TestHTTP_SetConnectionClose(t *testing.T) {
	t.Parallel()

	h := Get().SetConnectionClose()
	req, _, _ := h.test()
	if !req.Header.ConnectionClose() {
		t.Error("SetConnectionClose failed")
	}
}

func TestHTTP_SetRequestURI(t *testing.T) {
	t.Parallel()

	h := Get().SetRequestURI("https://ip.x2ox.com")
	req, _, _ := h.test()
	if req.URI().String() != "https://ip.x2ox.com/" {
		t.Errorf("SetRequestURI failed %s", req.URI().String())
	}
}

func TestHTTP_SetHost(t *testing.T) {
	t.Parallel()

	h := Get().SetHost("https://ip.x2ox.com")
	req, _, _ := h.test()
	if string(req.Host()) != "https://ip.x2ox.com" {
		t.Errorf("SetHost failed %s", string(req.Host()))
	}
}

func TestHTTPSetBodyWithEncoder(t *testing.T) {
	t.Parallel()

	h := Post().SetBodyWithEncoder(JSON(), struct {
		A string
	}{A: "A"})
	req, _, _ := h.test()
	if !bytes.Equal(req.Body(), []byte("{\"A\":\"A\"}\n")) {
		t.Error("err")
	}

	h = Post().SetBodyWithEncoder(nil, []byte("1"))
	req, _, _ = h.test()
	if !bytes.Equal(req.Body(), []byte("1")) {
		t.Error("err")
	}

	h = Post().SetBodyWithEncoder(nil, "1")
	req, _, _ = h.test()
	if !bytes.Equal(req.Body(), []byte("1")) {
		t.Error("err")
	}

	h = Post().SetBodyWithEncoder(nil, nil)
	req, _, _ = h.test()
	if !bytes.Equal(req.Body(), nil) {
		t.Error("err")
	}

	v := url.Values{}
	v.Add("A", "A")
	h = Post().SetBodyWithEncoder(nil, v)
	req, _, _ = h.test()
	if !bytes.Equal(req.Body(), []byte(v.Encode())) {
		t.Error("err")
	}
}
