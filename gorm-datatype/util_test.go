package datatype

import (
	"testing"
)

func TestScanLinearArray(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{``, "expected '{' at offset 0"},
		{`x`, "expected '{' at offset 0"},
		{`}`, "expected '{' at offset 0"},
		{`{`, "expected '}' at offset 1"},
		{`{{}`, "expected '}' at offset 3"},
		{`{}}`, "unexpected '}' at offset 2"},
		{`{,}`, "unexpected ',' at offset 1"},
		{`{,x}`, "unexpected ',' at offset 1"},
		{`{x,}`, "unexpected '}' at offset 3"},
		{`{x,{`, "unexpected '{' at offset 3"},
		{`{x},`, "unexpected ',' at offset 3"},
		{`{x}}`, "unexpected '}' at offset 3"},
		{`{{x}`, "expected '}' at offset 4"},
		{`{""x}`, "unexpected 'x' at offset 3"},
		{`{{a},{b,c}}`, "multidimensional arrays must have elements with matching dimensions"},
	} {
		_, err := scanLinearArray([]byte(tt.input), "")
		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}
	}
}
