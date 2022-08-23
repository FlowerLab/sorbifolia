package random

import (
	"math/rand"
	"time"
)

type MathRand struct {
	randBytes    []byte
	randBytesLen int
}

func init() {
	rand.Seed(time.Now().Unix())
}

func NewMathRand() RandString {
	return &MathRand{
		randBytes:    []byte(randBytes),
		randBytesLen: randBytesLen,
	}
}

func (r MathRand) RandString(length int) string {
	arr := make([]byte, length)
	for i := range arr {
		arr[i] = r.randBytes[uint64(rand.Uint32())*uint64(r.randBytesLen)>>32]
	}
	return string(arr)
}

func (r MathRand) SetRandBytes(data []byte) RandString {
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
