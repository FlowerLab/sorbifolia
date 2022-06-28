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
