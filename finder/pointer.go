package finder

func ToPtr[T any](x T) *T {
	return &x
}

func ToSlicePtr[T any](arr []T) []*T {
	out := make([]*T, len(arr))
	for _, v := range arr {
		out = append(out, ToPtr(v))
	}
	return out
}

func ToAnySlice[T any](arr []T) []any {
	out := make([]any, len(arr))
	for i, item := range arr {
		out[i] = item
	}
	return out
}
