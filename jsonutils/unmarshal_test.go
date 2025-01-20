package jsonutils

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestCompareToStd(t *testing.T) {
	tests := []string{
		`{}`,
		`{"a": 1}`,
		`{]`,
		`"abc"`,
		`5`,
		`{"a": 1} `,
		`{"a": 1} {}`,
		`{} bad data`,
		`{"a": 1} "hello"`,
		`[]`,
		`   {"x": {"t": [3,4,5]}}`,
	}

	for _, test := range tests {
		b := []byte(test)
		var ourV, stdV any
		ourErr := Unmarshal(b, &ourV)
		stdErr := json.Unmarshal(b, &stdV)
		if (ourErr == nil) != (stdErr == nil) {
			t.Errorf("Unmarshal(%q): our err = %#[2]v (%[2]T), std err = %#[3]v (%[3]T)", test, ourErr, stdErr)
		}

		if ourErr != nil && strings.HasPrefix(ourErr.Error(), "trailing garbage") {
			continue
		}

		if !reflect.DeepEqual(ourV, stdV) {
			t.Errorf("Unmarshal(%q): our val = %v, std val = %v", test, ourV, stdV)
		}
	}

	for _, test := range []string{
		"{\"abc\":\"abc\"}",
	} {
		b := []byte(test)
		ourV, stdV := make(map[int]string), make(map[int]string)
		ourErr := Unmarshal(b, &ourV)
		stdErr := json.Unmarshal(b, &stdV)
		if (ourErr == nil) != (stdErr == nil) {
			t.Errorf("Unmarshal(%q): our err = %#[2]v (%[2]T), std err = %#[3]v (%[3]T)", test, ourErr, stdErr)
		}

		if !reflect.DeepEqual(ourV, stdV) {
			t.Errorf("Unmarshal(%q): our val = %v, std val = %v", test, ourV, stdV)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	var m any
	j := []byte("1")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Unmarshal(j, &m)
	}
}

func BenchmarkStdUnmarshal(b *testing.B) {
	var m any
	j := []byte("1")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(j, &m)
	}
}
