package random

import (
	_ "runtime"
	_ "unsafe"
)

type fastRand struct{ *randString }

func Fast() Random {
	return &fastRand{newRandString()}
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

//go:linkname _fastRand runtime.fastrand
func _fastRand() uint32

//go:linkname _fastRand64 runtime.fastrand64
func _fastRand64() uint64
