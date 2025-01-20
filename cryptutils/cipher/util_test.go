package cipher

import (
	"bytes"
	"testing"
)

func TestSliceForAppend(t *testing.T) {
	t.Parallel()

	if head, tail := sliceForAppend(make([]byte, 10), 5); len(head) != cap(head) ||
		len(tail) != cap(tail) ||
		len(head) != 15 || len(tail) != 5. {
		t.Error("fail")
	}

	if head, tail := sliceForAppend(make([]byte, 0, 10), 5); len(head) != len(tail) ||
		cap(head) != cap(tail) ||
		len(head) != 5 || cap(head) != 10 {
		t.Error("fail")
	}
}

func TestGfnDouble(t *testing.T) {
	t.Parallel()

	if !bytes.Equal(gfnDouble([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}), []byte{
		2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32,
	}) {
		t.Error("fail")
	}

	defer func() {
		if err := recover(); err == nil {
			t.Error("fail")
		}
	}()
	gfnDouble([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
}
