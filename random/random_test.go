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

func TestPicks(t *testing.T) {
	arr := []int{123, 21, 21, 21, 3, 1233, 21, 321, 423, 4, 32, 43, 543, 5, 43}
	if len(Picks(arr, 3)) != 3 {
		t.Error("test fail")
	}
	if len(Picks(arr, 300)) != 300 {
		t.Error("test fail")
	}
	if Picks(arr, 0) != nil {
		t.Error("test fail")
	}
	if Picks([]int{}, 2) != nil {
		t.Error("test fail")
	}
}

func TestReverse(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6}
	Reverse(arr)
	if !isEqual(arr, []int{6, 5, 4, 3, 2, 1}) {
		t.Error("TestReverse failed: ")
	}
}

func TestRandom_Number(t *testing.T) {
	var arr = []Random{Safe(), Fast(), Math()}
	for _, v := range arr {
		v.Uint()
		v.Uint64()
		v.Uint32()
		v.Uint16()
		v.Uint8()
		v.Int()
		v.Int64()
		v.Int32()
		v.Int16()
		v.Int8()
		v.Uintn(12)
		v.Uint64n(32)
		v.Uint32n(21)
		v.Uint16n(12)
		v.Uint8n(33)
		v.Intn(11)
		v.Int64n(14)
		v.Int32n(16)
		v.Int16n(18)
		v.Int8n(10)

		v.Uintn(0)
		v.Uint64n(0)
		v.Uint32n(0)
		v.Uint16n(0)
		v.Uint8n(0)
		v.Intn(0)
		v.Int64n(0)
		v.Int32n(0)
		v.Int16n(0)
		v.Int8n(0)
	}
}
