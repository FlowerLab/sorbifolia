package random

import (
	"testing"
)

func TestShuffleAl(t *testing.T) {
	a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	b := make([]int, len(a))
	copy(b, a)

	Shuffle(a)

	if isEqual(a, b) {
		t.Errorf("Shuffle failed: %v | %v\n", a, b)
	}
}

func isEqual(item1, item2 []int) bool {
	if len(item1) != len(item2) {
		return false
	}
	if len(item1) == 0 || len(item2) == 0 {
		return len(item1) == 0 && len(item2) == 0
	}
	for i := 0; i < len(item1); i++ {
		if item1[i] != item2[i] {
			return false
		}
	}
	return true
}
