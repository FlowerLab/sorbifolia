package random

import (
	_ "runtime"
	_ "unsafe"
)

type FastRand struct {
	randBytes    []byte
	randBytesLen int
}

func NewFastRand() RandString {
	return &FastRand{
		randBytes:    []byte(randBytes),
		randBytesLen: randBytesLen,
	}
}

func (r FastRand) RandString(length int) string {
	arr := make([]byte, length)
	for i := range arr {
		arr[i] = r.randBytes[uint64(fastRand())*uint64(r.randBytesLen)>>32]
	}
	return string(arr)
}

func (r FastRand) SetRandBytes(data []byte) RandString {
	if len(data) > 256 {
		panic("data too long")
	}
	if hasRepeat(data) {
		panic("not repeatable")
	}

	r.randBytes = data
	r.randBytesLen = len(data)
	return r
}

//go:linkname fastRand runtime.fastrand
func fastRand() uint32

//go:linkname fastRand64 runtime.fastrand64
func fastRand64() uint64

func hasRepeat[T comparable](arr []T) bool {
	m := make(map[T]struct{})
	for _, v := range arr {
		if _, ok := m[v]; ok {
			return ok
		}
		m[v] = struct{}{}
	}
	return false
}
