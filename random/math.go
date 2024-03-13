package random

import (
	"math/rand/v2"
)

type mathRand struct{ *randString }

func Math() Random {
	return &mathRand{defaultRandString}
}

func (r mathRand) RandString(length int) string {
	arr := make([]int, length)
	for i := range arr {
		arr[i] = int(uint64(rand.Uint32()) * uint64(r.randBytesLen) >> 32)
	}
	return r.randString.RandString(arr)
}

func (r mathRand) SetRandBytes(data []byte) Random {
	return &mathRand{r.randString.SetRandBytes(data)}
}

func (r mathRand) Uint() uint     { return uint(rand.Uint64()) }
func (r mathRand) Uint64() uint64 { return rand.Uint64() }
func (r mathRand) Uint32() uint32 { return rand.Uint32() }
func (r mathRand) Uint16() uint16 { return uint16(rand.Uint32()) }
func (r mathRand) Uint8() uint8   { return uint8(rand.Uint32()) }

func (r mathRand) Int() int     { return toPtr[int](r.Uint()) }
func (r mathRand) Int64() int64 { return toPtr[int64](r.Uint64()) }
func (r mathRand) Int32() int32 { return toPtr[int32](r.Uint32()) }
func (r mathRand) Int16() int16 { return toPtr[int16](r.Uint16()) }
func (r mathRand) Int8() int8   { return toPtr[int8](r.Uint8()) }

func (r mathRand) Uintn(n uint) uint       { return remainder(n, r.Uint()) }
func (r mathRand) Uint64n(n uint64) uint64 { return remainder(n, r.Uint64()) }
func (r mathRand) Uint32n(n uint32) uint32 { return remainder(n, r.Uint32()) }
func (r mathRand) Uint16n(n uint16) uint16 { return remainder(n, r.Uint16()) }
func (r mathRand) Uint8n(n uint8) uint8    { return remainder(n, r.Uint8()) }

func (r mathRand) Intn(n int) int       { return remainder(n, r.Int()) }
func (r mathRand) Int64n(n int64) int64 { return remainder(n, r.Int64()) }
func (r mathRand) Int32n(n int32) int32 { return remainder(n, r.Int32()) }
func (r mathRand) Int16n(n int16) int16 { return remainder(n, r.Int16()) }
func (r mathRand) Int8n(n int8) int8    { return remainder(n, r.Int8()) }
