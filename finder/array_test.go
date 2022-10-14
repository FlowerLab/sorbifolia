package finder

import (
	"testing"
)

func TestIndexOf(t *testing.T) {
	t.Parallel()

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

func TestLastIndexOf(t *testing.T) {
	t.Parallel()

	if LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 2) != 4 {
		t.Error("LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 2) != 4")
	}
	if LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 6) != -1 {
		t.Error("LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 6) != -1")
	}
}

func TestFind(t *testing.T) {
	t.Parallel()

	if res, ok := Find([]string{"a", "b", "c", "d"}, func(i string) bool {
		return i == "b"
	}); !ok || res != "b" {
		t.Error("Find err")
	}

	if res, ok := Find([]string{"foobar"}, func(i string) bool {
		return i == "b"
	}); ok || res != "" {
		t.Error("Find err")
	}
}

func TestFindIndexOf(t *testing.T) {
	t.Parallel()

	if item, idx := FindIndexOf([]string{"a", "b", "c", "d", "b"}, func(i string) bool {
		return i == "b"
	}); idx != 1 || item != "b" {
		t.Error("Find index err")
	}

	if item, idx := FindIndexOf([]string{"foobar"}, func(i string) bool {
		return i == "b"
	}); idx != -1 || item != "" {
		t.Error("Find index err")
	}
}

func TestFindLastIndexOf(t *testing.T) {
	t.Parallel()

	if item, idx := FindLastIndexOf([]string{"a", "b", "c", "d", "b"}, func(i string) bool {
		return i == "b"
	}); idx != 4 || item != "b" {
		t.Error("Find index err")
	}

	if item, idx := FindLastIndexOf([]string{"foobar"}, func(i string) bool {
		return i == "b"
	}); idx != -1 || item != "" {
		t.Error("Find index err")
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	if !Contains([]string{"a", "b", "c", "d", "b"}, "b") {
		t.Error("Contains err")
	}
	if Contains([]string{"a", "b", "c", "d", "b"}, "e") {
		t.Error("Contains err")
	}
}

func TestContainsBy(t *testing.T) {
	t.Parallel()

	if !ContainsBy([]string{"a", "b", "c", "d", "b"}, func(i string) bool {
		return i == "b"
	}) {
		t.Error("ContainsBy err")
	}
	if ContainsBy([]string{"a", "b", "c", "d", "b"}, func(i string) bool {
		return i == "e"
	}) {
		t.Error("ContainsBy err")
	}
}
