package datatype

import (
	"testing"
)

func TestINetAddr(t *testing.T) {
	testdata := []testStruct{
		{&INetAddr{}, "1.1.1.1", "1.1.1.1", false},
		{&INetAddr{}, "0.0.0.0", "0.0.0.0", false},
		{&INetAddr{}, "1.1.1.", "", true},
		{&INetAddr{}, "1.1.1.256", "", true},
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

func TestINetPrefix(t *testing.T) {
	testdata := []testStruct{
		{&INetPrefix{}, "1.1.1.1/32", "1.1.1.1/32", false},
		{&INetPrefix{}, "0.0.0.0/32", "0.0.0.0/32", false},
		{&INetPrefix{}, "1.1.1.1/0", "1.1.1.1/0", false},
		{&INetPrefix{}, "1.1.1.256", "", true},
		{&INetPrefix{}, "1.1.1.256/-1", "", true},
		{&INetPrefix{}, "1.1.1.256/33", "", true},
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
