package datatype

import (
	"testing"
)

func TestArrayString(t *testing.T) {
	testdata := []testStruct{
		{&ArrayString{}, `{"1","2","3"}`, `{"1","2","3"}`, false},
		{&ArrayString{"1", "2", "3"}, nil, `{"1","2","3"}`, false},
		{&ArrayString{}, "{}", nil, false},
		{&ArrayString{}, map[string]string{}, nil, true},
		{&ArrayString{}, `{{"1"}}`, nil, true},
		{&ArrayString{}, nil, nil, false},

		{&ArrayString{}, ``, nil, true},
		{&ArrayString{}, `{`, nil, true},
		{&ArrayString{}, `{{a},{b}}`, nil, true},
		{&ArrayString{}, `{NULL}`, nil, true},
		{&ArrayString{}, `{a,NULL}`, nil, true},
		{&ArrayString{}, `{a,b,NULL}`, nil, true},
	}

	for _, v := range testdata {
		if v.data != nil {
			if err := v.itr.Scan(v.data); err != nil && !v.isErr {
				t.Error(err)
			}
		}

		val, err := v.itr.Value()
		if v.val != val {
			t.Error("is not a valid", v.val, val)
		}
		if err != nil && !v.isErr {
			t.Error(err)
		}
	}
}
