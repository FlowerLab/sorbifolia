package random

import (
	"crypto/rand"
	"unsafe"

	pn "go.x2ox.com/sorbifolia/pyrokinesis"
)

type safeRand struct{ *randString }

func Safe() Random {
	return &safeRand{defaultRandString}
}

func (r safeRand) RandString(length int) string {
	arr := make([]byte, length)
	arr1 := make([]int, length)
	_, _ = rand.Read(arr)
	for i := range arr {
		arr1[i] = int(uint64(arr[i]) * uint64(r.randBytesLen) >> 8)
	}
	return r.randString.RandString(arr1)
}

func (r safeRand) SetRandBytes(data []byte) Random {
	return &safeRand{r.randString.SetRandBytes(data)}
}

func (r safeRand) Uint() uint     { return pn.Bytes.ToUint(r.read(_intSize)) }
func (r safeRand) Uint64() uint64 { return pn.Bytes.ToUint64(r.read(8)) }
func (r safeRand) Uint32() uint32 { return pn.Bytes.ToUint32(r.read(4)) }
func (r safeRand) Uint16() uint16 { return pn.Bytes.ToUint16(r.read(2)) }
func (r safeRand) Uint8() uint8   { return pn.Bytes.ToUint8(r.read(1)) }

func (r safeRand) Int() int     { return pn.Bytes.ToInt(r.read(_intSize)) }
func (r safeRand) Int64() int64 { return pn.Bytes.ToInt64(r.read(8)) }
func (r safeRand) Int32() int32 { return pn.Bytes.ToInt32(r.read(4)) }
func (r safeRand) Int16() int16 { return pn.Bytes.ToInt16(r.read(2)) }
func (r safeRand) Int8() int8   { return pn.Bytes.ToInt8(r.read(1)) }

func (r safeRand) Uintn(n uint) uint       { return remainder(n, r.Uint()) }
func (r safeRand) Uint64n(n uint64) uint64 { return remainder(n, r.Uint64()) }
func (r safeRand) Uint32n(n uint32) uint32 { return remainder(n, r.Uint32()) }
func (r safeRand) Uint16n(n uint16) uint16 { return remainder(n, r.Uint16()) }
func (r safeRand) Uint8n(n uint8) uint8    { return remainder(n, r.Uint8()) }

func (r safeRand) Intn(n int) int       { return remainder(n, r.Int()) }
func (r safeRand) Int64n(n int64) int64 { return remainder(n, r.Int64()) }
func (r safeRand) Int32n(n int32) int32 { return remainder(n, r.Int32()) }
func (r safeRand) Int16n(n int16) int16 { return remainder(n, r.Int16()) }
func (r safeRand) Int8n(n int8) int8    { return remainder(n, r.Int8()) }

func (r safeRand) read(num int) []byte {
	arr := make([]byte, num)
	_, _ = rand.Read(arr)
	return arr
}

const _intSize = int(unsafe.Sizeof(int(0)))
