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
	r.AddRoute(GET, "/api/v1/data", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/user/:id", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/file/*file", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/go/:id/bind", []Handler[string]{func(*string) {}})

	t.Run("", func(t *testing.T) {
		defer func() { _ = recover() }()
		r.AddRoute(GET, "api/v1/data", []Handler[string]{func(*string) {}})
		t.Error("fail")
	})

	t.Run("", func(t *testing.T) {
		defer func() { _ = recover() }()
		r.AddRoute(GET, "/api/v1//data", []Handler[string]{func(*string) {}})
		t.Error("fail")
	})

	t.Run("", func(t *testing.T) {
		defer func() { _ = recover() }()
		r.AddRoute(GET, "/api/v1/data", nil)
		t.Error("fail")
	})

	t.Run("", func(t *testing.T) {
		defer func() { _ = recover() }()
		r.AddRoute(Method(10), "/api/v1/data", []Handler[string]{func(*string) {}})
		t.Error("fail")
	})
}

func TestRouter_Group(t *testing.T) {
	r := NewRouter[string]()
	g := r.Group()

	g.Group("/api/v1", func(*string) {})
	g.Use(func(*string) {})
	user := g.Group("/user")
	user.GET("/", func(*string) {})
	user.POST("/", func(*string) {})
	user.DELETE("/", func(*string) {})
	user.PATCH("/", func(*string) {})
	user.PUT("/", func(*string) {})
	user.OPTIONS("/", func(*string) {})
	user.HEAD("/", func(*string) {})
	user.CONNECT("/", func(*string) {})
	user.TRACE("/", func(*string) {})

	g.Any("/file/", func(*string) {})
}

func TestRouter_Find(t *testing.T) {
	r := NewRouter[string]()
	r.AddRoute(GET, "/api/v1/data", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v1/a/b/:id/d", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/user/:id", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/file/*file", []Handler[string]{func(*string) {}})
	r.AddRoute(GET, "/api/v2/go/:id/bind", []Handler[string]{func(*string) {}})
	r.Sort()

	ps := &Params{}
	if a := r.Find(GET, "/api/v1/a/b/c/d", ps); a == nil || (*ps)[0].Val != "c" {
		t.Error("fail")
	}
	if _, ok := ps.Get(":i"); ok {
		t.Error("fail")
	}
	if a := r.Find(GET, "/api/v1/a/b/c/e", ps); a != nil {
		t.Error("fail")
	}
	if a := r.Find(GET, "/api/v2/file/a/b/c", ps); a == nil || (*ps)[0].Val != "a/b/c" {
		t.Error("fail")
	}
	if val, ok := ps.Get("*file"); !ok || val != "a/b/c" {
		t.Error("fail")
	}
}
