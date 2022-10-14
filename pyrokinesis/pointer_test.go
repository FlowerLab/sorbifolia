package pyrokinesis

import (
	"testing"
)

func TestPointer(t *testing.T) {
	t.Parallel()

	i := 1
	ptr := Ptr(&i)
	ii := To[int](ptr)
	if *ii != i {
		t.Error("fail")
	}
}
