package httprouter

import (
	"testing"
)

func TestCheckDuplication(t *testing.T) {
	t.Parallel()

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
	checkDuplication(n)

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
	defer func() { _ = recover() }()
	checkDuplication(n)
	t.Error("duplication check failed")
}

func BenchmarkCheckDuplication(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		checkDuplication(n)
	}
}

func TestCheckNodeType(t *testing.T) {

	t.Parallel()

	t.Run("", func(t *testing.T) {
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
		checkNodeType(n)
	})

	t.Run("", func(t *testing.T) {
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
					{Type: NodeFixed},
					{Type: NodeFixed},
					{Type: NodeWild},
				}},
				{Type: NodeFixed},
				{Type: NodeWild},
			},
		}
		defer func() { _ = recover() }()
		checkNodeType(n)
		t.Error("fail")
	})

	t.Run("", func(t *testing.T) {
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
					{Type: NodeFixed},
					{Type: NodeWild},
					{Type: NodeWild},
				}},
				{Type: NodeFixed},
				{Type: NodeWild},
			},
		}
		defer func() { _ = recover() }()
		checkNodeType(n)
		t.Error("fail")
	})

	t.Run("", func(t *testing.T) {
		n := &Node[string]{
			Type: NodeStatic,
			ChildNode: []*Node[string]{
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeStatic},
				{Type: NodeWild, ChildNode: []*Node[string]{
					{Type: NodeStatic},
					{Type: NodeStatic},
					{Type: NodeStatic},
				}},
				{Type: NodeFixed},
			},
		}
		defer func() { _ = recover() }()
		checkNodeType(n)
		t.Error("fail")
	})
}

func TestSortName(t *testing.T) {
	t.Parallel()

	t.Run("", func(t *testing.T) {
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
		sortNode(n)
	})
	t.Run("", func(t *testing.T) {
		n := &Node[string]{}
		sortNode(n)
	})
}

func BenchmarkCheckNodeType(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		checkNodeType(n)
	}
}

func BenchmarkSortNode(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		sortNode(n)
	}
}
