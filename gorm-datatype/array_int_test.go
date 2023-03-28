package datatype

import (
	"testing"

	"gorm.io/gorm"
)

type TestArrayInt struct {
	gorm.Model

	Data Array[int] `json:"data"`
}

type testStruct struct {
	itr   gormItr
	data  any // nil string []byte
	val   any // nil string []byte
	isErr bool
}

func TestArray(t *testing.T) {
	testdata := []testStruct{
		{&Array[int]{}, "{1,2,3,4,5,6,7,8}", "{1,2,3,4,5,6,7,8}", false},
		{&Array[int]{1, 2, 3, 4, 5, 6, 7, 8}, nil, "{1,2,3,4,5,6,7,8}", false},
		{&Array[int]{}, "{}", nil, false},
		{&Array[int]{}, map[string]string{}, nil, true},
		{&Array[int]{}, `{{"1"}}`, nil, true},
		{&Array[int]{}, `{3.3,2.7}`, nil, true},
		{&Array[int]{}, nil, nil, false},
	}

	for _, v := range testdata {
		if v.data != nil {
			if err := v.itr.Scan(v.data); err != nil && !v.isErr {
				t.Error(err)
			}
		}

		val, err := v.itr.Value()
		if v.val != val {
			t.Error("is not a valid")
		}
		if err != nil && !v.isErr {
			t.Error(err)
		}
	}
}
