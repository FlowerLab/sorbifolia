package pool

import (
	"sync"
)

var pools = sync.Map{}

func Get[T any]() *T {
	if val, ok := pools.Load(*new(T)); ok {
		return val.(*sync.Pool).Get().(*T)
	}
	return new(T)
}

func Put[T any](t *T) {
	if t == nil {
		return
	}
	key := *new(T)

	val, ok := pools.Load(key)
	if !ok {
		val = &sync.Pool{New: func() interface{} { return new(T) }}
		pools.Store(key, val)
	}
	val.(*sync.Pool).Put(t)
}
