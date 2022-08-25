package httprouter

import (
	"net/http"
	"testing"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	server   *Server
	handlers Handlers[Context]
	index    int
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

type Server struct {
	IRouter[Context]
	router *Router[Context]
}

func NewServer() *Server {
	r := NewRouter[Context]()
	return &Server{r.Group(), r}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{Request: r, Writer: w}
	_ = c
}

func (s Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}

func testHTTPRouter(t *testing.T) {
	server := NewServer()
	server.GET("/", func(c *Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		_, _ = c.Writer.Write([]byte(`{"a":"a"}`))
	})

	_ = server.ListenAndServe(":8080")
}
