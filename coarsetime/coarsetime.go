package coarsetime

import (
	"sync/atomic"
	"time"
)

var (
	coarseTime atomic.Value
	frequency  = time.Millisecond * 100
)

func init() {
	t := time.Now().Truncate(frequency)
	coarseTime.Store(&t)
	go func() {
		for {
			time.Sleep(frequency)
			nt := time.Now().Truncate(frequency)
			coarseTime.Store(&nt)
		}
	}()
}

func CeilingTime() time.Time { return Now().Add(frequency) }
func FloorTime() time.Time   { return Now() }
