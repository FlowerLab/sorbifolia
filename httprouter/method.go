package httprouter

type MethodNode[T any] struct {
	Method Method
	Node   *Node[T]
}

type Method uint8

const (
	GET Method = iota
	HEAD
	POST
	PUT
	PATCH
	DELETE
	CONNECT
	OPTIONS
	TRACE
)

var methods = []Method{GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE}
