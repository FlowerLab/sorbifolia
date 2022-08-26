package httprouter

import (
	"net/http"
	"testing"
)

func BenchmarkOneRoute(B *testing.B) {
	server := NewServer()
	server.GET("/", func(c *Context) {})
	runRequest(B, server, "GET", "/ping")
}

func BenchmarkManyRoutesFist(B *testing.B) {
	router := NewServer()
	router.Any("/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkOneRouteJson(B *testing.B) {
	router := NewServer()
	router.GET("/hi", func(c *Context) {
		c.JSON("hello world")
	})
	runRequest(B, router, "GET", "/hi")
}

func Benchmark404(B *testing.B) {
	router := NewServer()
	router.Any("/something", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func runRequest(B *testing.B, r *Server, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}

func (m *mockWriter) Header() (h http.Header) {
	return m.headers
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockWriter) WriteHeader(int) {}
