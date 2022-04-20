package random

import (
	cr "crypto/rand"
)

type SafeRand struct {
	randBytes    []byte
	randBytesLen int
}

func NewSafeRand() *SafeRand {
	return &SafeRand{
		randBytes:    []byte(randBytes),
		randBytesLen: randBytesLen,
	}
}

func (r SafeRand) RandString(length int) string {
	arr := make([]byte, length)
	if _, err := cr.Read(arr); err != nil {
		return ""
	}
	for i := range arr {
		arr[i] = r.randBytes[int(arr[i])%r.randBytesLen]
	}
	return string(arr)
}

func (r SafeRand) SetRandBytes(data []byte) RandString {
	r.randBytes = data
	r.randBytesLen = len(data)
	return r
}
