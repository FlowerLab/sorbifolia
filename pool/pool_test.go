package pool

import (
	"sync"
	"testing"
)

func TestGet(t *testing.T) {
	s := Get[string]()
	if s == nil {
		t.Fail()
	}

	Put(s)
}

func TestPool(t *testing.T) {
	p := NewPool[string](nil, nil)
	a := p.Get()
	p.Put(a)
}

func TestPoolStore(t *testing.T) {
	Put[int](nil)
	if _, ok := pools.Load(0); ok {
		t.Error("fail")
	}

	defer func() {
		if err := recover(); err == nil {
			t.Error("fail")
		}
	}()

	p := sync.Pool{New: func() interface{} { return new(string) }} //nolint:all
	p.Put(new(string))                                             //nolint:all
	pools.Store(0, &p)                                             //nolint:all
	Get[int]()
}
