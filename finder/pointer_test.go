package finder

import (
	"testing"
)

func TestToPtr(t *testing.T) {
	t.Parallel()

	ptr := ToPtr(123)
	if *ptr != 123 {
		t.Error("fail")
	}
}

func TestToSlicePtr(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5}
	ptr := ToSlicePtr(arr)
	for i, v := range ptr {
		if *v != i+1 {
			t.Error("fail")
		}
	}
}

func TestToAnySlice(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4, 5}
	res := ToAnySlice(arr)
	for i, v := range res {
		val, ok := v.(int)
		if ok && val == i+1 {
			continue
		}
		t.Error("fail")
	}
}
