package pool

import (
	"testing"
)

func TestPool(t *testing.T) {
	t.Parallel()

	p := NewPool[string](nil, nil)
	a := p.Get()
	p.Put(a)
}
