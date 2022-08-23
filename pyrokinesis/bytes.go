package pyrokinesis

import (
	"unsafe"
)

func (_Bytes) ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
