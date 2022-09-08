package random

import (
	cr "crypto/rand"
)

type safeRand struct{ *randString }

func Safe() Random {
	return &safeRand{newRandString()}
}

func (r safeRand) RandString(length int) string {
	arr := make([]byte, length)
	arr1 := make([]int, length)
	_, _ = cr.Read(arr)
	for i := range arr {
		arr1[i] = int(uint64(arr[i]) * uint64(r.randBytesLen) >> 8)
	}
	return r.randString.RandString(arr1)
}

func (r safeRand) SetRandBytes(data []byte) Random {
	return &safeRand{r.randString.SetRandBytes(data)}
}
