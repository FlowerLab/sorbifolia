package pyrokinesis

import (
	"reflect"
	"unsafe"
)

func (_String) Copy(s string) string {
	return string(String.ToBytes(s))
}

func (_String) ToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}))
}
