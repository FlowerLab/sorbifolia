package datatype

import (
	"bytes"
	"testing"
)

func TestScanLinearArray(t *testing.T) {
	for _, tt := range []struct {
		input string
		dims  []int
		elems [][]byte
	}{
		{`{}`, nil, [][]byte{}},
		{`{NULL}`, []int{1}, [][]byte{nil}},
		{`{a}`, []int{1}, [][]byte{{'a'}}},
		{`{""}`, []int{1}, [][]byte{{}}},
		{`{","}`, []int{1}, [][]byte{{','}}},
		{`{",",","}`, []int{2}, [][]byte{{','}, {','}}},
		{`{"\"}"}`, []int{1}, [][]byte{{'"', '}'}}},
		{`{"\"","\""}`, []int{2}, [][]byte{{'"'}, {'"'}}},
	} {
		_, err := scanLinearArray([]byte(tt.input), "")
		if err != nil {
			t.Fatalf("Expected no error for %q, got %q", tt.input, err)
		}
	}
}

func TestScanLinearArrayErr(t *testing.T) {
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
		{`{{",",","}}`, "cannot convert ARRAY"},
	} {
		_, err := scanLinearArray([]byte(tt.input), "")
		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}
	}
}

func TestAppendArrayQuotedBytes(t *testing.T) {
	for _, tt := range []struct {
		input, output string
	}{
		{``, `""`},
		{`\`, `"\\"`},
		{`"`, `"\""`},
		{`{`, `"{"`},
		{`{asd\`, `"{asd\\"`},
	} {
		out := appendArrayQuotedBytes(nil, []byte(tt.input))
		if !bytes.Equal(out, []byte(tt.output)) {
			t.Errorf("error for %q", tt.input)
		}
	}
}
