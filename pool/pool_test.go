package pool

import (
	"testing"
)

func TestGet(t *testing.T) {
	s := Get[string]()
	if s == nil {
		t.Fail()
	}

	Put(s)
}
