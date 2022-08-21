package datatype

import (
	"testing"

	"gorm.io/gorm"
)

func TestArrayString_GormDataType(t *testing.T) {
	var i ArrayString
	if i.GormDataType() != "ArrayString" {
		t.Error("GormDataType Array")
	}
}

func TestArrayString_GormDBDataType(t *testing.T) {
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
