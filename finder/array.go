package finder

// IndexOf returns the index at which the first occurrence is found in an array or -1
func IndexOf[T comparable](arr []T, e T) int {
	for i, item := range arr {
		if item == e {
			return i
		}
	}

	return -1
}

// LastIndexOf returns the index at which the last occurrence is found in an array or -1
func LastIndexOf[T comparable](arr []T, e T) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == e {
			return i
		}
	}

	return -1
}

// Find search an e in a slice based on a function
func Find[T any](arr []T, is func(T) bool) (T, bool) {
	for _, v := range arr {
		if is(v) {
			return v, true
		}
	}
	return *new(T), false
}

// FindIndexOf searches an e in a slice based on a function and returns the index or -1
func FindIndexOf[T any](arr []T, is func(T) bool) (T, int) {
	for i, v := range arr {
		if is(v) {
			return v, i
		}
	}

	return *new(T), -1
}

// FindLastIndexOf searches last e in a slice based on a function and returns the index or -1
func FindLastIndexOf[T any](arr []T, is func(T) bool) (T, int) {
	for i := len(arr) - 1; i >= 0; i-- {
		if is(arr[i]) {
			return arr[i], i
		}
	}

	return *new(T), -1
}

// Contains returns true if an element is present in a collection.
func Contains[T comparable](arr []T, e T) bool {
	return IndexOf(arr, e) != -1
}

// ContainsBy returns true if function return true.
func ContainsBy[T any](arr []T, is func(T) bool) bool {
	_, ok := Find(arr, is)
	return ok
}
