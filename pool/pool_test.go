package pool

import (
	"testing"
)

func TestPool(t *testing.T) {
	p := NewPool[string](nil, nil)
	a := p.Get()
	p.Put(a)
}
