package datatype

import (
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func TestGormDataType(t *testing.T) {
	testdata := []struct {
		i   gormItr
		res string
	}{
		{&Array[int]{}, "Array"},
		{&ArrayString{}, "ArrayString"},
		{&JSON{}, "json"},
		{&INetAddr{}, "INetAddr"},
	}
	for _, v := range testdata {
		testGormDataType(t, v.i, v.res)
	}
}

func TestGormDBDataType(t *testing.T) {
	testdata := []struct {
		i   gormItr
		res string
	}{
		{&Array[int]{}, "text[]"},
		{&ArrayString{}, "text[]"},
		{&JSON{}, "JSONB"},
		{&INetAddr{}, "inet"},
	}
	for _, v := range testdata {
		testGormDBDataType(t, v.i, v.res)
	}
}

func testGormDataType(t *testing.T, i gormItr, res string) {
	if i.GormDataType() != res {
		t.Error("GormDataType no match")
	}
}

func testGormDBDataType(t *testing.T, i gormItr, res string) {
	if i.GormDBDataType(&gorm.DB{Config: &gorm.Config{Dialector: &_testPD{name: "postgres"}}}, nil) != res {
		t.Error("GormDBDataType no match", res)
	}
	if i.GormDBDataType(&gorm.DB{Config: &gorm.Config{Dialector: &_testPD{name: "mysql"}}}, nil) == res {
		t.Error("GormDBDataType no match", res)
	}
}

type _testPD struct {
	name string
}

func (t *_testPD) Name() string                                      { return t.name }
func (*_testPD) Initialize(_ *gorm.DB) error                         { panic("?") }
func (*_testPD) Migrator(_ *gorm.DB) gorm.Migrator                   { panic("?") }
func (*_testPD) DataTypeOf(_ *schema.Field) string                   { panic("?") }
func (*_testPD) DefaultValueOf(_ *schema.Field) clause.Expression    { panic("?") }
func (*_testPD) BindVarTo(_ clause.Writer, _ *gorm.Statement, _ any) { panic("?") }
func (*_testPD) QuoteTo(_ clause.Writer, _ string)                   { panic("?") }
func (*_testPD) Explain(_ string, _ ...any) string                   { panic("?") }
