package datatype

import (
	"testing"
)

func TestJSON(t *testing.T) {
	testdata := []testStruct{
		{&JSON{}, `{"a":"a"}`, `{"a":"a"}`, false},
		{&JSON{}, `{"a":"}`, nil, true},
	}

	for _, v := range testdata {
		if v.data != nil {
			if err := v.itr.Scan(v.data); err != nil && !v.isErr {
				t.Error(err)
			}
		}

		val, err := v.itr.Value()
		if v.val != val {
			t.Errorf("is not a valid ext: %s, act: %s", v.val, val)
		}
		if err != nil && !v.isErr {
			t.Error(err)
		}
	}
}
