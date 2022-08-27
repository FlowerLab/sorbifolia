package httprouter

import (
	"strings"
)

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

func NewMethod(method string) Method {
	switch strings.ToUpper(method) {
	case "GET":
		return GET
	case "HEAD":
		return HEAD
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "PATCH":
		return PATCH
	case "DELETE":
		return DELETE
	case "CONNECT":
		return CONNECT
	case "OPTIONS":
		return OPTIONS
	case "TRACE":
		return TRACE
	}

	return 255
}
