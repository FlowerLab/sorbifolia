package datatype

import (
	"testing"

	"gorm.io/gorm"
)

type TestArrayInt struct {
	gorm.Model

	Data Array[int] `json:"data"`
}

func TestArray_Int(t *testing.T) {
	initDB()
	if err := db.AutoMigrate(&TestArrayInt{}); err != nil {
		t.Error(err)
	}

	t1 := &TestArrayInt{}
	t1.Data = []int{
		123, 123, 123, 213, 23, 0, 32, 32, 32, 32, 32, 2, 3, 2332, 123, 23,
	}

	if err := db.Create(t1).Error; err != nil {
		t.Error(err)
	}
}

func TestArray_GormDataType(t *testing.T) {
	t.Parallel()

	var i Array[int]
	if i.GormDataType() != "Array" {
		t.Error("GormDataType Array")
	}
}

func TestArray_GormDBDataType(t *testing.T) {
	t.Parallel()

	var (
		i      Array[int]
		testDB = &gorm.DB{Config: &gorm.Config{Dialector: testDialector("postgres")}}
	)

	if i.GormDBDataType(testDB, nil) != "text[]" {
		t.Error("fail")
	}
	testDB = &gorm.DB{Config: &gorm.Config{Dialector: testDialector("mysql")}}
	if i.GormDBDataType(testDB, nil) != "" {
		t.Error("fail")
	}
}

func TestArray_Scan(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		str string
		arr []int
	}{
		{`{}`, nil},
		{`{1}`, []int{1}},
		{`{3,7}`, []int{3, 7}},
		{`{3,1,2}`, []int{3, 1, 2}},
	} {
		bytes := []byte(tt.str)
		a := Array[int]{}
		err := a.Scan(bytes)

		if err != nil {
			t.Errorf("expected no error for %q, got %v", bytes, err)
		}
	}

	for _, tt := range []struct {
		str string
		arr []int
	}{
		{`{}`, nil},
		{`{1}`, []int{1}},
		{`{3,7}`, []int{3, 7}},
		{`{3,1,2}`, []int{3, 1, 2}},
	} {
		a := Array[int]{}
		err := a.Scan(tt.str)

		if err != nil {
			t.Errorf("expected no error for %s, got %v", tt.str, err)
		}
	}

	a := Array[int]{}
	if err := a.Scan(nil); err != nil {
		t.Errorf("expected no error for nil, got %v", err)
	}

	if err := a.Scan(map[string]interface{}{}); err == nil {
		t.Error("fail")
	}
}

func TestArray_ScanErr(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		str string
		arr []int
	}{
		{`{{"1"}}`, []int{1}},
		{`{3.3,2.7}`, []int{3, 7}},
	} {
		bytes := []byte(tt.str)
		a := Array[int]{}
		err := a.Scan(bytes)

		if err == nil {
			t.Errorf("expected no error for %q, got %v", bytes, err)
		}
	}
}

func TestArray_Value(t *testing.T) {
	t.Parallel()

	var a Array[int] = nil
	if val, err := a.Value(); err != nil || val != nil {
		t.Error("expected")
	}

	a = []int{}
	if val, err := a.Value(); err != nil || val != "{}" {
		t.Error("expected")
	}
}
