package pyrokinesis

import (
	"testing"
)

func TestPointer(t *testing.T) {
	i := 1
	ptr := Ptr(&i)
	ii := To[int](ptr)
	if *ii != i {
		t.Error("fail")
	}
}
