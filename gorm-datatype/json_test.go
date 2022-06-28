package datatype

import (
	"encoding/json"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	if db, err = gorm.Open(
		postgres.Open("host=127.0.0.1 user=u123456789 password=u123456789 dbname=postgres port=5432 sslmode=disable"),
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
