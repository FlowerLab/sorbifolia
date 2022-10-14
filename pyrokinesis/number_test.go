package pyrokinesis

import (
	"testing"
)

func TestNumber_ToBytes(t *testing.T) {
	t.Parallel()

	Number[int]{}.ToBytes(1)
	Number[int64]{}.ToBytes(1)
	Number[int32]{}.ToBytes(1)
	Number[int16]{}.ToBytes(1)
	Number[int8]{}.ToBytes(1)

	Number[uint]{}.ToBytes(1)
	Number[uint64]{}.ToBytes(1)
	Number[uint32]{}.ToBytes(1)
	Number[uint16]{}.ToBytes(1)
	Number[uint8]{}.ToBytes(1)
}
