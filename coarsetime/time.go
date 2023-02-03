package coarsetime

import (
	"time"
)

func Since(t time.Time) time.Duration { return FloorTime().Sub(t) }
func Until(t time.Time) time.Duration { return t.Sub(FloorTime()) }
func Now() time.Time                  { return *coarseTime.Load().(*time.Time) }
func DateTime() string                { return Now().Format(time.DateTime) }

// Ptr return *time.Time, this cannot be changed, will affect the use of other goroutines
func Ptr() *time.Time { return coarseTime.Load().(*time.Time) }
