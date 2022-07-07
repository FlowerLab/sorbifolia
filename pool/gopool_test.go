package pool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

const runTimes uint32 = 1000

var sum uint32

func fn(i uint32) {
	atomic.AddUint32(&sum, i)
	fmt.Printf("run with %d\n", i)
}

func TestGoPool(t *testing.T) {
	var wg sync.WaitGroup
	pool := NewGoPool(10, func(i uint32) {
		fn(i)
		wg.Done()
	})
	for i := uint32(0); i < runTimes; i++ {
		wg.Add(1)
		pool.Invoke(i)
	}
	wg.Wait()
	t.Logf("finish all tasks, result is %d", sum)
}
