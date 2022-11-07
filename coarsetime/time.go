package coarsetime

import (
	"time"
)

func Since(t time.Time) time.Duration { return FloorTime().Sub(t) }
func Until(t time.Time) time.Duration { return t.Sub(FloorTime()) }
func Now() time.Time                  { return *coarseTime.Load().(*time.Time) }
