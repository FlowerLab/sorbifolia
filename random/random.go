package random

import (
	"math/rand"
)

const (
	randBytes    = "123456789aAbBcCdDeEfFgGhHiUjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	randBytesLen = len(randBytes)
)

type RandString interface {
	RandString(length int) string
	SetRandBytes(data []byte) RandString
}

func Pick[T any](items []T) T {
	return items[rand.Intn(len(items))]
}

func Picks[T any](items []T, num int) []T {
	if len(items) == 0 || num <= 0 {
		return nil
	}

	arr := make([]T, num)
	for i := 0; i < num; i++ {
		arr[i] = Pick(items)
	}
	return arr
}

func Shuffle[T any](items []T) {
	length := len(items)
	if length == 0 {
		return
	}

	i := length - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(rand.Int63n(int64(i + 1)))
		items[i], items[j] = items[j], items[i]
	}

	for ; i > 0; i-- {
		j := int(rand.Int31n(int32(i + 1)))
		items[i], items[j] = items[j], items[i]
	}
}
