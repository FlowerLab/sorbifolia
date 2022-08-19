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

// CeilingTime
//
// Deprecated: As of Go 1.17, Performance is lower than time.Now().
func CeilingTime() time.Time {
	tp := coarseTime.Load().(*time.Time)
	return (*tp).Add(frequency)
}

func FloorTime() time.Time {
	tp := coarseTime.Load().(*time.Time)
	return *tp
}
