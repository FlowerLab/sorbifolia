package datastructures

import (
	"fmt"
	"strings"
)

type Set[T comparable] struct {
	items map[T]struct{}
}

func (set *Set[T]) Add(items ...T) {
	for _, item := range items {
		set.items[item] = struct{}{}
	}
}

func (set *Set[T]) Remove(items ...T) {
	for _, item := range items {
		delete(set.items, item)
	}
}

func (set *Set[T]) Exists(item T) bool {
	_, ok := set.items[item]
	return ok
}

func (set *Set[T]) Pop() T {
	for item := range set.items {
		delete(set.items, item)
		return item
	}
	return *new(T)
}

func (set *Set[T]) Has(items ...T) bool {
	if len(items) == 0 {
		return false
	}

	for _, item := range items {
		if _, has := set.items[item]; !has {
			return false
		}
	}
	return true
}

func (set *Set[T]) List() []T {
	arr := make([]T, 0, len(set.items))
	for item := range set.items {
		arr = append(arr, item)
	}
	return arr
}

func (set *Set[T]) Len() int64    { return int64(len(set.items)) }
func (set *Set[T]) IsEmpty() bool { return set.Len() == 0 }

// IsEqual test whether s and t are the same in size and have the same items.
func (set *Set[T]) IsEqual(t *Set[T]) bool {
	if set.Len() != t.Len() {
		return false
	}

	equal := true
	t.Each(func(item T) bool {
		_, equal = set.items[item]
		return equal
	})

	return equal
}

// IsSubset tests whether t is a subset of set.
func (set *Set[T]) IsSubset(t *Set[T]) (subset bool) {
	subset = true

	t.Each(func(item T) bool {
		_, subset = set.items[item]
		return subset
	})

	return
}

// IsSuperset tests whether t is a superset of set.
func (set *Set[T]) IsSuperset(t *Set[T]) bool { return t.IsSubset(set) }

// Copy returns a new Set with a copy of set.
func (set *Set[T]) Copy() *Set[T] {
	u := New[T]()
	u.Add(set.List()...)
	return u
}

// String returns a string representation of set
func (set *Set[T]) String() string {
	arr := set.List()
	t := make([]string, 0, len(arr))
	for _, item := range arr {
		t = append(t, fmt.Sprintf("%v", item))
	}

	return fmt.Sprintf("[%s]", strings.Join(t, ", "))
}

func (set *Set[T]) Merge(t *Set[T]) {
	t.Each(func(item T) bool {
		set.items[item] = struct{}{}
		return true
	})
}

func (set *Set[T]) Separate(t *Set[T]) { set.Remove(t.List()...) }

func (set *Set[T]) Each(f func(item T) bool) {
	for item := range set.items {
		if !f(item) {
			break
		}
	}
}

func (set *Set[T]) Clear() { set.items = map[T]struct{}{} }

func (set *Set[T]) Union(sets ...*Set[T]) *Set[T] {
	s := set.Copy()

	for _, v := range sets {
		v.Each(func(item T) bool {
			s.Add(item)
			return true
		})
	}
	return s
}

func (set *Set[T]) Diff(sets ...*Set[T]) *Set[T] {
	s := set.Copy()
	for _, v := range sets {
		s.Separate(v)
	}
	return s
}

func (set *Set[T]) Intersection(sets ...*Set[T]) *Set[T] {
	all := set.Union(sets...)
	result := New[T]()

	all.Each(func(item T) bool {
		if has := set.Has(item); !has {
			return true
		}
		for _, s := range sets {
			if !s.Has(item) {
				return true
			}
		}
		result.Add(item)
		return true
	})
	return result
}

func (set *Set[T]) SymmetricDiff(s *Set[T]) *Set[T] {
	return set.Diff(s).Union(s.Diff(set))
}

func New[T comparable](items ...T) *Set[T] {
	set := &Set[T]{items: make(map[T]struct{})}
	set.Add(items...)
	return set
}
