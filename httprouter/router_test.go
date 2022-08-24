package httprouter

import (
	"testing"
)

func TestNewRouter(t *testing.T) {
	r := NewRouter[string]()
	if len(r.method) != 9 {
		t.Error("fail")
	}
	for _, v := range r.method {
		if v == nil {
			t.Error("fail")
		}
	}
	_ = r.Method(GET).Type
	_ = r.Method(HEAD).Type
	_ = r.Method(POST).Type
	_ = r.Method(PUT).Type
	_ = r.Method(PATCH).Type
	_ = r.Method(DELETE).Type
	_ = r.Method(CONNECT).Type
	_ = r.Method(OPTIONS).Type
	_ = r.Method(TRACE).Type

	defer func() { _ = recover() }()
	_ = r.Method(Method(10)).Type

	t.Error("fail")
}

func TestRouter_AddRoute(t *testing.T) {
	r := NewRouter[string]()
	r.AddRoute(GET, "/api/v1/data", []HandlerFunc[string]{func(*string) {}})
	r.AddRoute(GET, "/", []HandlerFunc[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/user/:id", []HandlerFunc[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/file/*file", []HandlerFunc[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/go/:id/bind", []HandlerFunc[string]{func(*string) {}})

	t.Run("", func(t *testing.T) {
		defer func() { _ = recover() }()
		r.AddRoute(GET, "api/v1/data", []HandlerFunc[string]{func(*string) {}})
		t.Error("fail")
	})

	t.Run("", func(t *testing.T) {
		defer func() { _ = recover() }()
		r.AddRoute(GET, "/api/v1//data", []HandlerFunc[string]{func(*string) {}})
		t.Error("fail")
	})

	t.Run("", func(t *testing.T) {
		defer func() { _ = recover() }()
		r.AddRoute(GET, "/api/v1/data", nil)
		t.Error("fail")
	})

	t.Run("", func(t *testing.T) {
		defer func() { _ = recover() }()
		r.AddRoute(Method(10), "/api/v1/data", []HandlerFunc[string]{func(*string) {}})
		t.Error("fail")
	})
}
