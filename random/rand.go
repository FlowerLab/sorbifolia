package random

import (
	"math/rand"
)

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
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
}

func Reverse[T any](collection []T) {
	length := len(collection)

	for i := 0; i < length/2; i = i + 1 {
		j := length - 1 - i
		collection[i], collection[j] = collection[j], collection[i]
	}
}
