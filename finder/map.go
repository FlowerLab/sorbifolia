package finder

// Keys creates an array of the map keys.
func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}

// Values creates an array of the map values.
func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))

	for _, v := range m {
		result = append(result, v)
	}

	return result
}

// PickBy returns same map type filtered by given predicate.
func PickBy[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	r := map[K]V{}
	for k, v := range m {
		if f(k, v) {
			r[k] = v
		}
	}
	return r
}

// PickByKeys returns same map type filtered by given keys.
func PickByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V {
	r := map[K]V{}
	for k, v := range m {
		if Contains(keys, k) {
			r[k] = v
		}
	}
	return r
}

// PickByValues returns same map type filtered by given values.
func PickByValues[K comparable, V comparable](m map[K]V, values []V) map[K]V {
	r := map[K]V{}
	for k, v := range m {
		if Contains(values, v) {
			r[k] = v
		}
	}
	return r
}

// OmitBy returns same map type filtered by given predicate.
func OmitBy[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	r := map[K]V{}
	for k, v := range m {
		if !f(k, v) {
			r[k] = v
		}
	}
	return r
}

// OmitByKeys returns same map type filtered by given keys.
func OmitByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V {
	r := map[K]V{}
	for k, v := range m {
		if !Contains(keys, k) {
			r[k] = v
		}
	}
	return r
}

// OmitByValues returns same map type filtered by given values.
func OmitByValues[K comparable, V comparable](in map[K]V, values []V) map[K]V {
	r := map[K]V{}
	for k, v := range in {
		if !Contains(values, v) {
			r[k] = v
		}
	}
	return r
}

// Invert creates a map composed of the inverted keys and values. If map
// contains duplicate values, subsequent values overwrite property assignments
// of previous values.
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	out := map[V]K{}

	for k, v := range m {
		out[v] = k
	}

	return out
}

// Assign merges multiple maps from left to right.
func Assign[K comparable, V any](maps ...map[K]V) map[K]V {
	out := map[K]V{}

	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}

	return out
}

// MapKeys manipulates a map keys and transforms it to a map of another type.
func MapKeys[K comparable, V any, R comparable](m map[K]V, f func(V, K) R) map[R]V {
	result := map[R]V{}

	for k, v := range m {
		result[f(v, k)] = v
	}

	return result
}

// MapValues manipulates a map values and transforms it to a map of another type.
func MapValues[K comparable, V any, R any](m map[K]V, f func(V, K) R) map[K]R {
	result := map[K]R{}

	for k, v := range m {
		result[k] = f(v, k)
	}

	return result
}
