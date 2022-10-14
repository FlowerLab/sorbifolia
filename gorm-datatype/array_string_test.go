package datatype

import (
	"testing"

	"gorm.io/gorm"
)

func TestArrayString_GormDataType(t *testing.T) {
	t.Parallel()

	var i ArrayString
	if i.GormDataType() != "ArrayString" {
		t.Error("GormDataType Array")
	}
}

func TestArrayString_GormDBDataType(t *testing.T) {
	t.Parallel()

	var (
		i      ArrayString
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

func TestArrayString_Scan(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		str string
		arr []string
	}{
		{`{}`, nil},
		{`{t}`, []string{"t"}},
		{`{f,1}`, []string{"f", "1"}},
		{`{"a\\b","c d",","}`, []string{"a\\b", "c d", ","}},
	} {
		bytes := []byte(tt.str)
		a := ArrayString{}
		err := a.Scan(bytes)

		if err != nil {
			t.Errorf("expected no error for %q, got %v", bytes, err)
		}
	}

	for _, tt := range []struct {
		str string
		arr []string
	}{
		{`{}`, nil},
		{`{t}`, []string{"t"}},
		{`{f,1}`, []string{"f", "1"}},
		{`{"a\\b","c d",","}`, []string{"a\\b", "c d", ","}},
	} {
		a := ArrayString{}
		err := a.Scan(tt.str)

		if err != nil {
			t.Errorf("expected no error for %s, got %v", tt.str, err)
		}
	}

	a := ArrayString{}
	if err := a.Scan(nil); err != nil {
		t.Errorf("expected no error for nil, got %v", err)
	}

	if err := a.Scan(map[string]interface{}{}); err == nil {
		t.Error("fail")
	}
}

func TestArrayString_ScanErr(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		str string
		err string
	}{
		{``, "unable to parse array"},
		{`{`, "unable to parse array"},
		{`{{a},{b}}`, "cannot convert ARRAY[2][1] to StringArray"},
		{`{NULL}`, "parsing array element index 0: cannot convert nil to string"},
		{`{a,NULL}`, "parsing array element index 1: cannot convert nil to string"},
		{`{a,b,NULL}`, "parsing array element index 2: cannot convert nil to string"},
	} {
		bytes := []byte(tt.str)
		a := ArrayString{}
		err := a.Scan(bytes)

		if err == nil {
			t.Errorf("expected no error for %q, got %v", bytes, err)
		}
	}
}

func TestArrayString_Value(t *testing.T) {
	t.Parallel()

	var a ArrayString = nil
	if val, err := a.Value(); err != nil || val != nil {
		t.Error("expected")
	}

	a = []string{}
	if val, err := a.Value(); err != nil || val != "{}" {
		t.Error("expected")
	}

	a = []string{"a", "b", "c", "d", "e", "f"}
	if val, err := a.Value(); err != nil || val != "{\"a\",\"b\",\"c\",\"d\",\"e\",\"f\"}" {
		t.Error("expected")
	}
}
