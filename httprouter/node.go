package httprouter

import (
	"strings"
)

type Node[T any] struct {
	Path      string
	Type      NodeType
	Wild      bool
	Handler   Handlers[T]
	ChildNode []*Node[T]
}

func (n *Node[T]) AddRoute(path string, handlers Handlers[T]) {
	n.addNodeRoute(path[1:], handlers)
}

func (n *Node[T]) addNodeRoute(nodePath string, handlers Handlers[T]) {
	var (
		paths = strings.SplitN(nodePath, "/", 2)
		node  = n.getOrAddChildNode(paths[0])
	)

	if len(paths) == 1 { // end of traversal
		node.Handler = append(node.Handler, handlers...)
		return
	}
	node.addNodeRoute(paths[1], handlers)
}

func (n *Node[T]) getOrAddChildNode(path string) *Node[T] {
	nodeType := NodeStatic
	switch {
	case len(path) == 0:
	case path[0] == ':':
		nodeType = NodeFixed
	case path[0] == '*':
		nodeType = NodeWild
	}

	for _, v := range n.ChildNode {
		if v.Path == path {
			return v
		}
	}

	node := &Node[T]{Path: path, Type: nodeType}
	n.ChildNode = append(n.ChildNode, node)
	return node
}

type NodeType uint8

const (
	NodeStatic NodeType = iota
	NodeFixed
	NodeWild
)
