package pool

import (
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
