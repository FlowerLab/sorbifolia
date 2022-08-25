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

func (r *Router[T]) AddRoute(method Method, path string, handlers Handlers[T]) {
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

func (r *Router[T]) Group() IRouter[T] {
	return &Group[T]{route: r}
}

func (r *Router[T]) Sort() {
	for _, v := range methods {
		if r.Method(v).Type == NodeStatic &&
			len(r.Method(v).Handler) == 0 &&
			len(r.Method(v).ChildNode) == 0 {
			continue
		}
		checkNodeType(r.method[v])
		checkDuplication(r.method[v])
		sortNode(r.method[v])
	}
}

func (r *Router[T]) Find(method Method, path string) (Handlers[T], Params) {
	node, params := r.findNode(r.Method(method), strings.Split(path, "/"))
	if node == nil {
		return nil, nil
	}
	return node.Handler, params
}

func (r *Router[T]) findNode(node *Node[T], path []string) (*Node[T], Params) {
	var params Params

	switch node.Type {
	case NodeWild:
		return node, Params{{node.Path, strings.Join(path, "/")}}
	case NodeFixed:
		params = append(params, Param{node.Path, path[0]})
	case NodeStatic:
		if path[0] != node.Path {
			return nil, nil
		}
	}

	if len(path) == 1 {
		return node, params
	}
	for _, v := range node.ChildNode {
		if val, p := r.findNode(v, path[1:]); val != nil {
			params = append(params, p...)
			return val, params
		}
	}
	return nil, nil
}
