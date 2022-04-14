package random

type RandString interface {
	RandString(length int) string
	SetRandBytes(data []byte) RandString
}
