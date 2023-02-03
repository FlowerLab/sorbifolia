package pyrokinesis

import (
	"unsafe"
)

func (_Bytes) ToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func (_Bytes) ToInt(b []byte) int     { return toNumber[int](b) }
func (_Bytes) ToInt64(b []byte) int64 { return toNumber[int64](b) }
func (_Bytes) ToInt32(b []byte) int32 { return toNumber[int32](b) }
func (_Bytes) ToInt16(b []byte) int16 { return toNumber[int16](b) }
func (_Bytes) ToInt8(b []byte) int8   { return toNumber[int8](b) }

func (_Bytes) ToUint(b []byte) uint     { return toNumber[uint](b) }
func (_Bytes) ToUint64(b []byte) uint64 { return toNumber[uint64](b) }
func (_Bytes) ToUint32(b []byte) uint32 { return toNumber[uint32](b) }
func (_Bytes) ToUint16(b []byte) uint16 { return toNumber[uint16](b) }
func (_Bytes) ToUint8(b []byte) uint8   { return toNumber[uint8](b) }

func toNumber[T int | int64 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8](b []byte) T {
	var (
		val  T
		size = len(b)
	)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = b[i]
	}
	return val
}
