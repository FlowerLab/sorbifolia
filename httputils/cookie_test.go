package httputils

import (
	"testing"
)

func TestHTTP_SetCookie(t *testing.T) {
	h := Connect().SetCookie("foo", "bar")
	req, _, _ := h.test()
	if string(req.Header.Cookie("foo")) != "bar" {
		t.Error("SetCookie err")
	}
}

func TestHTTP_DelCookie(t *testing.T) {
	h := Options().SetCookie("foo", "bar").DelCookie("foo")

	req, _, _ := h.test()
	if string(req.Header.Cookie("foo")) != "" {
		t.Error("DelCookie err")
	}
}

func TestHTTP_DelAllCookies(t *testing.T) {
	h := Delete().SetCookie("foo", "bar").SetCookie("123", "321").DelAllCookies()

	req, _, _ := h.test()
	if string(req.Header.Cookie("foo")) != "" || string(req.Header.Cookie("123")) != "" {
		t.Error("DelAllCookies err")
	}
}
