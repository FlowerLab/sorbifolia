package random

const (
	randBytes    = "123456789aAbBcCdDeEfFgGhHiUjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	randBytesLen = len(randBytes)
)

type Random interface {
	RandString(length int) string
	SetRandBytes(data []byte) Random
}
