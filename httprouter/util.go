package httprouter

import (
	"sort"
)

type priorityNode[T any] struct {
	priority int
	node     *Node[T]
}

type priorityNodes[T any] []priorityNode[T]

func (n priorityNodes[T]) Len() int           { return len(n) }
func (n priorityNodes[T]) Less(i, j int) bool { return n[i].priority < n[j].priority }
func (n priorityNodes[T]) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

func sortNode[T any](n *Node[T]) {
	if len(n.ChildNode) <= 1 {
		return
	}

	arr := make(priorityNodes[T], len(n.ChildNode))
	for i, v := range n.ChildNode {
		arr[i].node = v

		switch v.Type {
		case NodeWild:
			arr[i].priority = -2
		case NodeFixed:
			arr[i].priority = -1
		default:
			arr[i].priority = getPriority(v)
		}
	}
	sort.Sort(sort.Reverse(arr))

	for i, v := range arr {
		n.ChildNode[i] = v.node
		if len(v.node.ChildNode) > 1 {
			sortNode(v.node)
		}
	}
}

func getPriority[T any](n *Node[T]) int {
	if n.Type == NodeWild || n.Type == NodeFixed {
		return 0
	}

	priority := 1
	for _, v := range n.ChildNode {
		priority += getPriority(v)
	}
	return priority
}

func checkNodeType[T any](n *Node[T]) {
	if n == nil || len(n.ChildNode) == 0 {
		return
	}

	var hasFixed, hasWild = false, false

	for _, v := range n.ChildNode {
		switch v.Type {
		case NodeWild:
			if hasWild {
				panic("only one wild can exist in the same path " + v.Path)
			}
			hasWild = true
			if len(v.ChildNode) != 0 {
				panic("wild cannot have follow-up paths")
			}
		case NodeFixed:
			if hasFixed {
				panic("only one fixed can exist in the same path " + v.Path)
			}
			hasFixed = true
		}

		checkNodeType(v)
	}

}

func checkDuplication[T any](n *Node[T]) {
	if n == nil || len(n.ChildNode) == 0 {
		return
	}

	set := make(map[string]struct{}, len(n.ChildNode))
	for _, v := range n.ChildNode {
		if _, ok := set[v.Path]; ok {
			panic("duplication of path: " + v.Path)
		}
		checkDuplication(v)
		set[v.Path] = struct{}{}
	}
}
