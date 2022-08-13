package httputils

import (
	"testing"
)

func TestHTTP_Add(t *testing.T) {
	h := Post().Add(nil)
	if len(h.fn) != 2 {
		t.Error("Add err")
	}
}
