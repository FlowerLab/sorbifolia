package httprouter

import (
	"testing"
)

func TestCheckDuplication(t *testing.T) {
	n := &Node[string]{
		ChildNode: []*Node[string]{
			{Path: "", ChildNode: []*Node[string]{
				{Path: "123", ChildNode: []*Node[string]{
					{Path: "312"},
				}},
				{Path: "", ChildNode: []*Node[string]{
					{Path: "123"},
				}},
				{Path: "321", ChildNode: []*Node[string]{
					{Path: "123"},
				}},
			}},
			{Path: "1", ChildNode: []*Node[string]{
				{Path: "a"},
			}},
			{Path: "2", ChildNode: []*Node[string]{
				{Path: "a"},
			}},
			{Path: "3", ChildNode: []*Node[string]{
				{Path: "a"},
			}},
		},
	}
	if checkDuplication(n) {
		t.Error("duplication check failed")
	}

	n = &Node[string]{
		ChildNode: []*Node[string]{
			{Path: "", ChildNode: []*Node[string]{
				{Path: "123", ChildNode: []*Node[string]{
					{Path: "312"},
				}},
				{Path: "", ChildNode: []*Node[string]{
					{Path: "123"},
				}},
				{Path: "", ChildNode: []*Node[string]{
					{Path: "123"},
				}},
			}},
			{Path: "1", ChildNode: []*Node[string]{
				{Path: "a"},
			}},
			{Path: "2", ChildNode: []*Node[string]{
				{Path: "a"},
			}},
			{Path: "3", ChildNode: []*Node[string]{
				{Path: "a"},
			}},
		},
	}
	if !checkDuplication(n) {
		t.Error("duplication check failed")
	}
}

func checkDuplication[T any](n *Node[T]) bool {
	if n == nil || len(n.ChildNode) == 0 {
		return false
	}

	set := make(map[string]struct{}, len(n.ChildNode))
	for _, v := range n.ChildNode {
		if _, ok := set[v.Path]; ok {
			return true
		}
		set[v.Path] = struct{}{}
	}
	for _, v := range n.ChildNode {
		if checkDuplication(v) {
			return true
		}
	}

	return false
}

func TestCheckNodeType(t *testing.T) {
	n := &Node[string]{
		Type: NodeStatic,
		ChildNode: []*Node[string]{
			{Type: NodeStatic},
			{Type: NodeStatic},
			{Type: NodeStatic},
			{Type: NodeStatic},
			{Type: NodeStatic, ChildNode: []*Node[string]{
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeFixed},
				{Type: NodeWild},
			}},
			{Type: NodeFixed},
			{Type: NodeWild},
		},
	}
	if checkNodeType(n) {
		t.Error("duplication NodeType check failed")
	}

	n = &Node[string]{
		Type: NodeStatic,
		ChildNode: []*Node[string]{
			{Type: NodeStatic},
			{Type: NodeStatic},
			{Type: NodeStatic},
			{Type: NodeStatic},
			{Type: NodeStatic, ChildNode: []*Node[string]{
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeFixed},
				{Type: NodeFixed},
				{Type: NodeWild},
			}},
			{Type: NodeFixed},
			{Type: NodeWild},
		},
	}
	if !checkNodeType(n) {
		t.Error("duplication NodeType check failed")
	}
}

func checkNodeType[T any](n *Node[T]) bool {
	if n == nil || len(n.ChildNode) == 0 {
		return false
	}

	var (
		hasFixed, hasWild = false, false
	)

	for _, v := range n.ChildNode {
		switch v.Type {
		case NodeWild:
			if hasWild {
				return true
			}
			hasWild = true
			if len(v.ChildNode) != 0 {
				panic("wild cannot have follow-up paths")
			}
		case NodeFixed:
			if hasFixed {
				return true
			}
			hasFixed = true
		}
	}
	for _, v := range n.ChildNode {
		if checkNodeType(v) {
			return true
		}
	}

	return false
}
