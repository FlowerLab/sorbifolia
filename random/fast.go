package random

import (
	_ "runtime"
	"unsafe"
	_ "unsafe"
)

type fastRand struct{ *randString }

func Fast() Random {
	return &fastRand{defaultRandString}
}

func (r fastRand) RandString(length int) string {
	arr := make([]int, length)
	for i := range arr {
		arr[i] = int(uint64(_fastRand()) * uint64(r.randString.randBytesLen) >> 32)
	}
	return r.randString.RandString(arr)
}

func (r fastRand) SetRandBytes(data []byte) Random {
	return &fastRand{r.randString.SetRandBytes(data)}
}

func (r fastRand) Uint() uint     { return _fastRandUint() }
func (r fastRand) Uint64() uint64 { return _fastRand64() }
func (r fastRand) Uint32() uint32 { return _fastRand() }
func (r fastRand) Uint16() uint16 { return uint16(_fastRand()) }
func (r fastRand) Uint8() uint8   { return uint8(_fastRand()) }

func (r fastRand) Int() int     { return toPtr[int](r.Uint()) }
func (r fastRand) Int64() int64 { return toPtr[int64](r.Uint64()) }
func (r fastRand) Int32() int32 { return toPtr[int32](r.Uint32()) }
func (r fastRand) Int16() int16 { return toPtr[int16](r.Uint16()) }
func (r fastRand) Int8() int8   { return toPtr[int8](r.Uint8()) }

func (r fastRand) Uintn(n uint) uint       { return remainder(n, r.Uint()) }
func (r fastRand) Uint64n(n uint64) uint64 { return remainder(n, r.Uint64()) }
func (r fastRand) Uint32n(n uint32) uint32 { return remainder(n, r.Uint32()) }
func (r fastRand) Uint16n(n uint16) uint16 { return remainder(n, r.Uint16()) }
func (r fastRand) Uint8n(n uint8) uint8    { return remainder(n, r.Uint8()) }

func (r fastRand) Intn(n int) int       { return remainder(n, r.Int()) }
func (r fastRand) Int64n(n int64) int64 { return remainder(n, r.Int64()) }
func (r fastRand) Int32n(n int32) int32 { return remainder(n, r.Int32()) }
func (r fastRand) Int16n(n int16) int16 { return remainder(n, r.Int16()) }
func (r fastRand) Int8n(n int8) int8    { return remainder(n, r.Int8()) }

func toPtr[T int | int64 | int32 | int16 | int8, N uint | uint64 | uint32 | uint16 | uint8](n N) T {
	return *(*T)(unsafe.Pointer(&n))
}

func remainder[T uint | uint64 | uint32 | uint16 | uint8 | int | int64 | int32 | int16 | int8](n, r T) T {
	if n == 0 {
		return r
	}
	return r % n
}

//go:linkname _fastRand runtime.fastrand
func _fastRand() uint32

//go:linkname _fastRand64 runtime.fastrand64
func _fastRand64() uint64

//go:linkname _fastRandUint runtime.fastrandu
func _fastRandUint() uint
