package httprouter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
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

var (
	// size                        = unsafe.Sizeof(&Node[string]{}) // and other overhead
	maxStackSize int32 = 1000 // max stack size 1<<20    100 for ci
	maxCount     int32 = 0    // recursive times
)

// When the Node's depth is 1048683 or 1064946 ,ths stack is overflow.
// Due to err of stack overflow,I have to use a file to record.
// The other overhead is maxCount*size
func TestCheckDuplicationDeep(t *testing.T) {
	maxCount = 0
	filepath, err := os.Getwd()
	if err != nil {
		panic("获取目录失败")
	}
	filepath = fmt.Sprintf("%s%s", filepath, "\\log.txt")
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		panic("文件打开失败")
	}
	defer func() { _ = recover() }()
	defer file.Close()
	writer := bufio.NewWriter(file)

	n := new(Node[string])
	tmp := n
	ch := make(chan int32)
	go func() {
		for maxCount <= maxStackSize {
			checkDuplication(tmp)
			ch <- maxCount
			n.ChildNode = []*Node[string]{
				{Path: ""},
			}
			n = n.ChildNode[0]
			atomic.AddInt32(&maxCount, 1)
		}
	}()
	for {
		count := <-ch
		if count >= maxStackSize {
			break
		}
		_, _ = writer.WriteString("递归的深度: " + strconv.Itoa(int(count)) + "\n")
		_ = writer.Flush()
	}
	close(ch)
}

func TestCheckNodeTypeDeep(t *testing.T) {
	maxCount = 0
	filepath, err := os.Getwd()
	if err != nil {
		panic("获取目录失败")
	}
	filepath = fmt.Sprintf("%s%s", filepath, "\\log.txt")
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		panic("文件打开失败")
	}
	defer func() { _ = recover() }()
	defer file.Close()
	writer := bufio.NewWriter(file)

	ch := make(chan int32)

	n := new(Node[string])
	tmp := n
	go func() {
		for maxCount <= maxStackSize {
			checkNodeType(tmp)
			ch <- maxCount
			n.ChildNode = []*Node[string]{
				{Path: ""},
			}
			n = n.ChildNode[0]
			atomic.AddInt32(&maxCount, 1)
		}
	}()
	for {
		count := <-ch
		if count >= maxStackSize {
			break
		}
		_, _ = writer.WriteString("递归的深度: " + strconv.Itoa(int(count)) + "\n")
		_ = writer.Flush()
	}
	close(ch)
}

func TestSortNodeDeep(t *testing.T) {
	maxCount = 0
	filepath, err := os.Getwd()
	if err != nil {
		panic("获取目录失败")
	}
	filepath = fmt.Sprintf("%s%s", filepath, "\\log.txt")
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		panic("文件打开失败")
	}
	defer func() { _ = recover() }()
	defer file.Close()
	writer := bufio.NewWriter(file)

	n := new(Node[string])
	tmp := n
	ch := make(chan int32)
	go func() {
		for maxCount <= maxStackSize {
			sortNode(tmp)
			n.ChildNode = []*Node[string]{
				{Path: ""},
			}
			n = n.ChildNode[0]
			ch <- maxCount
			atomic.AddInt32(&maxCount, 1)
		}
	}()
	for {
		count := <-ch
		if count >= maxStackSize {
			break
		}
		_, _ = writer.WriteString("递归的深度: " + strconv.Itoa(int(count)) + "\n")
		_ = writer.Flush()
	}
	close(ch)
}
