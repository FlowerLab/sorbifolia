package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type mathRand struct{ *randString }

func Math() Random {
	return &mathRand{newRandString()}
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
