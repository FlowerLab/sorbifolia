package random

import (
	cr "crypto/rand"
)

type SafeRand struct {
	randBytes    []byte
	randBytesLen int
}

func NewSafeRand() RandString {
	return &SafeRand{
		randBytes:    []byte(randBytes),
		randBytesLen: randBytesLen,
	}
}

func (r SafeRand) RandString(length int) string {
	arr := make([]byte, length)
	_, _ = cr.Read(arr)
	for i := range arr {
		arr[i] = r.randBytes[int(arr[i])%r.randBytesLen]
	}
	return string(arr)
}

func (r SafeRand) SetRandBytes(data []byte) RandString {
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
