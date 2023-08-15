package random

var defaultRandString = newRandString()

type randString struct {
	randBytes    []byte
	randBytesLen int
}

func newRandString() *randString {
	return &randString{
		randBytes:    []byte(randBytes),
		randBytesLen: randBytesLen,
	}
}

func (r randString) RandString(rn []int) string {
	arr := make([]byte, len(rn))
	for i := range arr {
		arr[i] = r.randBytes[rn[i]]
	}
	return string(arr)
}

func (r randString) SetRandBytes(data []byte) *randString {
	if len(data) > 256 {
		panic("data too long")
	}
	if hasRepeat(data) {
		panic("not repeatable")
	}
	r.randBytesLen = len(data)
	r.randBytes = data
	return &r
}

func hasRepeat[T comparable](arr []T) bool {
	m := make(map[T]struct{})
	for _, v := range arr {
		if _, ok := m[v]; ok {
			return ok
		}
		m[v] = struct{}{}
	}
	return false
}
