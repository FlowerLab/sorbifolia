package datastructures

type ListSingly[T any] struct {
	Data T
	next *ListSingly[T]
}

func (l *ListSingly[T]) Tail() *ListSingly[T]      { return l.next }
func (l *ListSingly[T]) Empty() bool               { return l == nil }
func (l *ListSingly[T]) Add(head T) *ListSingly[T] { return &ListSingly[T]{head, l} }

func (l *ListSingly[T]) Length() uint {
	curr, length := l, uint(0)
	for !curr.Empty() {
		curr, length = curr.Tail(), length+1
	}
	return length
}

func (l *ListSingly[T]) Insert(val T, pos uint) (*ListSingly[T], error) {
	if pos == 0 {
		return l.Add(val), nil
	}
	nl, err := l.next.Insert(val, pos-1)
	if err != nil {
		return nil, err
	}
	return nl.Add(l.Data), nil
}

func (l *ListSingly[T]) Get(pos uint) (T, bool) {
	if pos == 0 {
		return l.Data, true
	}
	return l.next.Get(pos - 1)
}

func (l *ListSingly[T]) Remove(pos uint) (*ListSingly[T], error) {
	if pos == 0 {
		return l.Tail(), nil
	}

	nl, err := l.next.Remove(pos - 1)
	if err != nil {
		return nil, err
	}
	return &ListSingly[T]{l.Data, nl}, nil
}

func (l *ListSingly[T]) Find(fn func(T) bool) (T, bool) {
	if fn(l.Data) {
		return l.Data, true
	}
	return l.next.Find(fn)
}

func (l *ListSingly[T]) FindIndex(fn func(T) bool) int {
	curr, idx := l, 0
	for !curr.Empty() {
		if fn(curr.Data) {
			return idx
		}
		curr, idx = curr.Tail(), idx+1
	}
	return -1
}

func (l *ListSingly[T]) Map(f func(T) T) []T { return append(l.next.Map(f), f(l.Data)) }
