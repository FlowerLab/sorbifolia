package random

const (
	randBytes    = "0123456789aAbBcCdDeEfFgGhHiUjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	randBytesLen = len(randBytes)
)

type Random interface {
	SetRandBytes(data []byte) Random
	RandString(length int) string

	Uint() uint
	Uint64() uint64
	Uint32() uint32
	Uint16() uint16
	Uint8() uint8

	Int() int
	Int64() int64
	Int32() int32
	Int16() int16
	Int8() int8

	Uintn(n uint) uint
	Uint64n(n uint64) uint64
	Uint32n(n uint32) uint32
	Uint16n(n uint16) uint16
	Uint8n(n uint8) uint8

	Intn(n int) int
	Int64n(n int64) int64
	Int32n(n int32) int32
	Int16n(n int16) int16
	Int8n(n int8) int8
}
