package reflectype

import (
	"reflect"
	"time"
)

var (
	Time         = reflect.TypeFor[time.Time]()
	TimeDuration = reflect.TypeFor[time.Duration]()
)
