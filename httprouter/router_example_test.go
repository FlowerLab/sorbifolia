package httprouter

import (
	"encoding/json"
	"net/http"
	"testing"
)

type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	server   *Server //nolint:all
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

func (c *Context) Status(code int) *Context { c.Writer.WriteHeader(code); return c }
func (c *Context) JSON(data interface{})    { _ = json.NewEncoder(c.Writer).Encode(data) }

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

//nolint:all
func testHTTPRouter(t *testing.T) {
	server := NewServer()
	server.GET("/", func(c *Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		_, _ = c.Writer.Write([]byte(`{"a":"a"}`))
	})

	_ = server.ListenAndServe(":8080")
}
