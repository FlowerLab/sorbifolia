package util

func ToNonNegativeInt(b []byte) (n int) {
	if len(b) == 0 {
		return -1
	}
	for _, val := range b {
		if val > '9' || val < '0' {
			return -1
		}
		n = n*10 + int(val-'0')
	}
	return
}
