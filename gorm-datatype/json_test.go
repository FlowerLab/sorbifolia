package datatype

import (
	"bytes"
	"encoding/json"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func initDB() {
	var err error
	if db, err = gorm.Open(
		postgres.Open("host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"),
	); err != nil {
		panic(err)
	}
}

type Test1 struct {
	gorm.Model

	Data JSON `json:"data"`
}

func TestJSON(t *testing.T) {
	initDB()
	if err := db.AutoMigrate(&Test1{}); err != nil {
		t.Error(err)
	}
	d := map[string]interface{}{
		"data": map[string]interface{}{
			"data": "ASd",
			"d":    123,
		},
	}
	t1 := &Test1{}
	bts, _ := json.Marshal(d)
	_ = json.Unmarshal(bts, &t1)
	if err := db.Create(t1).Error; err != nil {
		t.Error(err)
	}
}

func TestJSON_Value(t *testing.T) {
	var i JSON
	if val, err := i.Value(); err != nil || val != nil {
		t.Error("fail")
	}

	i = []byte(`{"a":"a"}`)
	if val, _ := i.Value(); val != `{"a":"a"}` {
		t.Error("fail")
	}
}

func TestJSON_Scan(t *testing.T) {
	var i JSON

	if err := i.Scan(nil); err != nil {
		t.Error(err)
	}
	if err := i.Scan([]byte(`{"a":"a"}`)); err != nil {
		t.Error(err)
	}
	if err := i.Scan(`{"a":"a"}`); err != nil {
		t.Error(err)
	}
	if err := i.Scan(map[string]string{"a": "a"}); err == nil {
		t.Error(err)
	}
}

func TestJSON_MarshalJSON(t *testing.T) {
	var i JSON
	if err := i.UnmarshalJSON([]byte(`{"a":"a"}`)); err != nil {
		t.Error(err)
	}
	if bts, err := i.MarshalJSON(); err != nil || !bytes.Equal(bts, []byte("{\"a\":\"a\"}")) {
		t.Error(err)
	}
}

func TestJSON_GormDataType(t *testing.T) {
	var i JSON
	if i.GormDataType() != "json" {
		t.Error("GormDataType json")
	}
}

func TestJSON_GormDBDataType(t *testing.T) {
	var (
		i      JSON
		testDB = &gorm.DB{Config: &gorm.Config{Dialector: testDialector("postgres")}}
	)

	if i.GormDBDataType(testDB, nil) != "JSONB" {
		t.Error("fail")
	}
	testDB = &gorm.DB{Config: &gorm.Config{Dialector: testDialector("mysql")}}
	if i.GormDBDataType(testDB, nil) != "" {
		t.Error("fail")
	}
}

type testDialector string

func (t testDialector) Name() string {
	return string(t)
}

func (t testDialector) Initialize(d *gorm.DB) error {
	panic("implement me")
}

func (t testDialector) Migrator(db *gorm.DB) gorm.Migrator {
	panic("implement me")
}

func (t testDialector) DataTypeOf(field *schema.Field) string {
	panic("implement me")
}

func (t testDialector) DefaultValueOf(field *schema.Field) clause.Expression {
	panic("implement me")
}

func (t testDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	panic("implement me")
}

func (t testDialector) QuoteTo(writer clause.Writer, s string) {
	panic("implement me")
}

func (t testDialector) Explain(sql string, vars ...interface{}) string {
	panic("implement me")
}
