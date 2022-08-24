package httprouter

import (
	"strings"
)

func NewRouter[T any]() *Router[T] {
	return &Router[T]{method: [9]*Node[T]{{}, {}, {}, {}, {}, {}, {}, {}, {}}}
}

type Router[T any] struct {
	method [9]*Node[T]
}

func (r *Router[T]) Method(method Method) *Node[T] { return r.method[method] }

func (r *Router[T]) AddRoute(method Method, path string, handlers HandlersChain[T]) {
	if path[0] != '/' {
		panic("path must begin with '/'")
	}
	if strings.Contains(path, "//") {
		panic("path cannot contain consecutive '/'")
	}
	if len(handlers) <= 0 {
		panic("there must be at least one handler")
	}

	r.Method(method).AddRoute(path, handlers)
}

func (r *Router[T]) Sort() error {

	return nil
}
