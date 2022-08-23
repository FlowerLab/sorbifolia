package pyrokinesis

import (
	"unsafe"
)

func To[T any](data unsafe.Pointer) *T  { return (*T)(data) }
func Ptr[T any](data *T) unsafe.Pointer { return unsafe.Pointer(data) }
