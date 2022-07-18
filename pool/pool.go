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

type Pool[T any] struct {
	pools *sync.Pool
	put   func(*T)
}

func (p *Pool[T]) Get() *T  { return p.pools.Get().(*T) }
func (p *Pool[T]) Put(t *T) { p.put(t); p.pools.Put(t) }

func NewPool[T any](get func() *T, put func(*T)) *Pool[T] {
	if get == nil {
		get = func() *T { return new(T) }
	}
	if put == nil {
		put = func(t *T) {}
	}

	return &Pool[T]{
		pools: &sync.Pool{New: func() any { return get() }},
		put:   put,
	}
}
