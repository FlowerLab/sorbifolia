package reflectype

import (
	"reflect"
)

var (
	Byte = reflect.TypeFor[byte]()
	Rune = reflect.TypeFor[rune]()
)
