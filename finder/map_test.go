package finder

import (
	"sort"
	"strconv"
	"testing"
)

func TestKeys(t *testing.T) {
	keys := Keys(map[string]int{"foo": 1, "bar": 2})
	sort.Strings(keys)

	if keys[0] != "bar" {
		t.Error("fail")
	}
}

func TestValues(t *testing.T) {
	values := Values(map[string]int{"foo": 1, "bar": 2})
	sort.Ints(values)

	if values[0] != 1 {
		t.Error("fail")
	}
}

func TestPickBy(t *testing.T) {
	val := PickBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
		return value > 2
	})
	if len(val) != 1 || Keys(val)[0] != "baz" {
		t.Error("fail")
	}
}

func TestPickByKeys(t *testing.T) {
	val := PickByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})

	if len(val) != 2 {
		t.Error("fail")
	}
}

func TestPickByValues(t *testing.T) {
	val := PickByValues(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})

	if len(val) != 2 {
		t.Error("fail")
	}
}

func TestOmitBy(t *testing.T) {
	val := OmitBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
		return value%2 == 1
	})

	if len(val) != 1 || Keys(val)[0] != "bar" {
		t.Error("fail")
	}
}

func TestOmitByKeys(t *testing.T) {
	val := OmitByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})

	if len(val) != 1 || Keys(val)[0] != "bar" {
		t.Error("fail")
	}
}

func TestOmitByValues(t *testing.T) {
	val := OmitByValues(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})

	if len(val) != 1 || Keys(val)[0] != "bar" {
		t.Error("fail")
	}
}

func TestInvert(t *testing.T) {
	v1 := Invert(map[string]int{"a": 1, "b": 2})
	v2 := Invert(map[string]int{"a": 1, "b": 2, "c": 1})

	if len(v1) != 2 || len(v2) != 2 {
		t.Error("fail")
	}
}

func TestAssign(t *testing.T) {
	v := Assign(map[string]int{"a": 1, "b": 2}, map[string]int{"b": 3, "c": 4})

	if len(v) != 3 {
		t.Error("fail")
	}
}

func TestMapKeys(t *testing.T) {
	v1 := MapKeys(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(_ int, _ int) string {
		return "a"
	})
	v2 := MapKeys(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(_ int, v int) string {
		return strconv.FormatInt(int64(v), 10)
	})

	if len(v1) != 1 || len(v2) != 4 {
		t.Error("fail")
	}
}

func TestMapValues(t *testing.T) {
	v1 := MapValues(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
		return "Hello"
	})
	v2 := MapValues(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
		return strconv.FormatInt(int64(x), 10)
	})

	if len(v1) != 4 || len(v2) != 4 {
		t.Error("fail")
	}
}
