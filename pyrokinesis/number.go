package pyrokinesis

import (
	"unsafe"
)

type Number[T int | int64 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8] struct {
}

func (Number[T]) ToBytes(num T) []byte {
	var (
		size = int(unsafe.Sizeof(num))
		arr  = make([]byte, size)
	)

	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}
