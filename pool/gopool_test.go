//go:build !race

package pool

import (
	"sync"
	"testing"
)

func TestGoPool(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	gp := NewGoPool(30, func(_ int) {
		wg.Done()
	})

	for i := 0; i < 10000; i++ {
		sum := i
		wg.Add(1)
		gp.Invoke(sum)
	}

	wg.Wait()
}
