package buffer

type Buffer interface {
	Bytes() []byte
	Len() int
	Cap() int
	Reset()
}
