package jsonutils

import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []string{
		`{"a":0}  {"a":1}  {"a":2} `,
		`{}  {"a":0}  {} `,
	}

	type test struct {
		A int `json:"a"`
	}

	for _, v := range tests {
		r := strings.NewReader(v)

		arr, err := Decode[test](r)
		if err != nil {
			t.Error(v, arr, err)
		}
	}
}
