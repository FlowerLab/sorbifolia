package reflectype

import (
	"reflect"
	"time"
)

var (
	Time         = reflect.TypeOf((*time.Time)(nil)).Elem()
	TimeDuration = reflect.TypeOf((*time.Duration)(nil)).Elem()
)
