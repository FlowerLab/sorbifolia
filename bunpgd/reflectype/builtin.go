package reflectype

import (
	"reflect"
)

var (
	Byte = reflect.TypeOf((*byte)(nil)).Elem()
	Rune = reflect.TypeOf((*rune)(nil)).Elem()
)
