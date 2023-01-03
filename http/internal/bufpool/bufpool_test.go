package bufpool

import (
	"fmt"
	"testing"
)

func TestAd(t *testing.T) {
	t.Log(-1 >> 1)
	t.Log(0 >> 1)
	t.Log(1230 >> 1)
	t.Log(64 >> 1)
}

func TestSF(t *testing.T) {
	arr := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}
	start, end := 3, 6
	arr = append(arr[:start], arr[end:]...)
	fmt.Println(arr)
}

// 	s1 := b.B[end:]
//	b.B = b.B[:start]
//	b.B = append(b.B, s1...)
