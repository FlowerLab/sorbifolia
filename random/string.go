package random

import (
	cr "crypto/rand"
	"math/rand"
	"time"
)

type (
	SafeRand struct {
		randBytes    []byte
		randBytesLen int
	}
	MathRand struct {
		randBytes    []byte
		randBytesLen int
	}
)

func NewSafeRand() RandString {
	return &SafeRand{
		randBytes:    []byte(randBytes),
		randBytesLen: randBytesLen,
	}
}
func NewMathRand() RandString {
	mr := &MathRand{
		randBytes:    []byte(randBytes),
		randBytesLen: randBytesLen,
	}
	rand.Seed(time.Now().UnixNano())
	return mr
}

const (
	randBytes    = "123456789aAbBcCdDeEfFgGhHiUjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	randBytesLen = len(randBytes)
)

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

func (r MathRand) RandString(length int) string {
	arr := make([]byte, length)
	for i := range arr {
		arr[i] = r.randBytes[rand.Intn(r.randBytesLen)]
	}
	return string(arr)
}

func (r MathRand) SetRandBytes(data []byte) RandString {
	r.randBytes = data
	r.randBytesLen = len(data)
	return r
}
