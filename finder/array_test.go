package finder

import (
	"testing"
)

func TestIndexOf(t *testing.T) {
	if IndexOf([]int{0, 1, 2, 1, 2, 3}, 2) != 2 {
		t.Error("IndexOf([]int{0, 1, 2, 1, 2, 3}, 2) != 2")
	}

	if IndexOf([]int{0, 1, 2, 1, 2, 3}, 4) != -1 {
		t.Error("IndexOf([]int{0, 1, 2, 1, 2, 3}, 3) != -1")
	}

	if IndexOf([]int{}, 4) != -1 {
		t.Error("IndexOf([]int{}, 4) != -1")
	}

	if IndexOf(nil, 4) != -1 {
		t.Error("IndexOf(nil, 3) != -1")
	}
}
