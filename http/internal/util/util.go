package util

func ToNonNegativeInt64(b []byte) (n int64) {
	for _, val := range b {
		if val > '9' || val < '0' {
			return -1
		}
		n = n*10 + int64(val-'0')
	}
	return
}
