package puresqlite

import (
	"database/sql"
	"fmt"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"modernc.org/sqlite"
)

func TestDialector(t *testing.T) {
	const (
		CustomDriverName = "my_custom_driver"
		InMemoryDSN      = "file:testdatabase?mode=memory&cache=shared"
	)

	sql.Register(CustomDriverName, &sqlite.Driver{})

	rows := []struct {
		description  string
		dialector    *Dialector
		openSuccess  bool
		query        string
		querySuccess bool
	}{
		{
			description: "Default driver",
			dialector: &Dialector{
				DSN: InMemoryDSN,
			},
			openSuccess:  true,
			query:        "SELECT 1",
			querySuccess: true,
		},
		{
			description: "Explicit default driver",
			dialector: &Dialector{
				DriverName: DriverName,
				DSN:        InMemoryDSN,
			},
			openSuccess:  true,
			query:        "SELECT 1",
			querySuccess: true,
		},
		{
			description: "Bad driver",
			dialector: &Dialector{
				DriverName: "not-a-real-driver",
				DSN:        InMemoryDSN,
			},
			openSuccess: false,
		},
		{
			description: "Custom driver",
			dialector: &Dialector{
				DriverName: CustomDriverName,
				DSN:        InMemoryDSN,
			},
			openSuccess:  true,
			query:        "SELECT 1",
			querySuccess: true,
		},
	}

	for rowIndex, row := range rows {
		t.Run(fmt.Sprintf("%d/%s", rowIndex, row.description), func(t *testing.T) {
			db, err := gorm.Open(row.dialector, &gorm.Config{})
			if !row.openSuccess {
				if err == nil {
					t.Errorf("Expected Open to fail.")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected Open to succeed; got error: %v", err)
			}
			if db == nil {
				t.Errorf("Expected db to be non-nil.")
			}
			if row.query != "" {
				err = db.Exec(row.query).Error
				if !row.querySuccess {
					if err == nil {
						t.Errorf("Expected query to fail.")
					}
					return
				}

				if err != nil {
					t.Errorf("Expected query to succeed; got error: %v", err)
				}
			}
		})
	}
}

func TestDialector_ORM(t *testing.T) {
	d := Open("file:testdb?mode=memory&cache=shared")
	t.Log(d.Name())
	db, err := gorm.Open(d, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.Exec("PRAGMA foreign_keys = 1")

	if db, err = gorm.Open(d, &gorm.Config{}); err != nil {
		t.Error(err)
	}

	db.Exec("PRAGMA foreign_keys = 1")

	type TestTable struct {
		gorm.Model

		Name string `json:"name"`
		Info string `json:"info"`
	}

	if err = db.AutoMigrate(&TestTable{}); err != nil {
		t.Error(err)
	}

	type UserTable struct {
		Username string `json:"username" gorm:"primarykey"`
		Info     string `json:"info"`
	}

	type DataTable struct {
		ID   uint    `json:"id" gorm:"primarykey;autoIncrement"`
		Info string  `json:"info"`
		F    float32 `json:"f"`
		B    []byte  `json:"b" gorm:"type:bytes"`
	}

	if err = db.AutoMigrate(&TestTable{}, &UserTable{}, &DataTable{}); err != nil {
		t.Error(err)
	}
	if err = db.Migrator().DropTable(&TestTable{}); err != nil {
		t.Error(err)
	}
	if err = db.Migrator().DropTable(&TestTable{}); err != nil {
		t.Error(err)
	}

	if err = db.Migrator().DropTable("SQLITE_MASTER"); err == nil {
		t.Fail()
	}

	type AMTable struct {
		ID uint   `json:"id" gorm:"primarykey;autoIncrement"`
		B  []byte `json:"b" gorm:"type:aaa"`
	}

	if err = db.AutoMigrate(&AMTable{}); err != nil {
		t.Error(err)
	}

	db.Create(&UserTable{
		Username: "a",
		Info:     "a",
	})
	db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "username"}}, DoNothing: true}).
		Create(&UserTable{
			Username: "a",
			Info:     "a",
		})

	var tut UserTable
	db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&tut, "username = ?", "a")
}
