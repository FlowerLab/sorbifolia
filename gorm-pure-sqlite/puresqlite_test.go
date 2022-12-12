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

func TestDialector_PreConn(t *testing.T) {
	d := Open("file:testdb1?mode=memory&cache=shared")
	db, err := gorm.Open(d, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.Exec("PRAGMA foreign_keys = 1")
	d.(*Dialector).Conn = db.ConnPool
	if err = d.Initialize(db); err != nil {
		t.Error(err)
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

	type TestTable struct {
		gorm.Model

		Name string `json:"name"`
		Info string `json:"info"`
	}

	if err = db.AutoMigrate(&TestTable{}); err != nil {
		t.Error(err)
	}
	_, _ = db.Migrator().GetTables()

	type UserTable struct {
		Username string `json:"username" gorm:"primarykey"`
		Info     string `json:"info"`
	}

	type DataTable struct {
		ID   uint    `json:"id" gorm:"primarykey;autoIncrement;default:1"`
		Info string  `json:"info" gorm:"default:1sdf"`
		BB   bool    `json:"bb"`
		F    float32 `json:"f"`
		B    []byte  `json:"b" gorm:"type:bytes"`
	}
	type ATTable struct {
		ID  uint `json:"id" gorm:"primarykey"`
		Num uint `json:"num" gorm:"autoIncrement"`
	}
	if err = db.AutoMigrate(&ATTable{}); err == nil {
		t.Fail()
	}
	if err = db.AutoMigrate(&TestTable{}, &UserTable{}, &DataTable{}); err != nil {
		t.Error(err)
	}
	{
		dt := &DataTable{BB: false}
		db.Create(dt)
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
	db.Clauses(testLocking{Strength: "UPDATE"}).First(&tut, "username = ?", "a")
	db.Clauses(testInsert{}).Create(&UserTable{Username: "c", Info: "c"})
	{
		var arr []UserTable
		db.Offset(1).Find(&arr)
	}

	_ = db.Transaction(func(tx *gorm.DB) error {
		tx.SavePoint("a1")
		tx.Create(&UserTable{Username: "a1", Info: "a"})
		tx.SavePoint("a2")
		tx.Create(&UserTable{Username: "a2", Info: "a"})

		if err = tx.Create(&UserTable{Username: "a1", Info: "a"}).Error; err != nil {
			tx.RollbackTo("a1")
		}

		return nil
	})
}

type testLocking struct {
	Strength string
	Table    clause.Table
	Options  string
}

// Name where clause name
func (locking testLocking) Name() string                      { return "FOR" }
func (locking testLocking) Build(_ clause.Builder)            {}
func (locking testLocking) MergeClause(clause *clause.Clause) { clause.Expression = locking }

type testInsert struct {
	Table    clause.Table
	Modifier string
}

func (insert testInsert) Name() string { return "INSERT" }

func (insert testInsert) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		_, _ = stmt.WriteString("INSERT ")
		if insert.Modifier != "" {
			_, _ = stmt.WriteString(insert.Modifier)
			_ = stmt.WriteByte(' ')
		}

		_, _ = stmt.WriteString("INTO ")
		if insert.Table.Name == "" {
			stmt.WriteQuoted(stmt.Table)
		} else {
			stmt.WriteQuoted(insert.Table)
		}
	}
}

// MergeClause merge insert clause
func (insert testInsert) MergeClause(clause *clause.Clause) {
	if v, ok := clause.Expression.(testInsert); ok {
		if insert.Modifier == "" {
			insert.Modifier = v.Modifier
		}
		if insert.Table.Name == "" {
			insert.Table = v.Table
		}
	}
	clause.Expression = insert
}

type TestAC1 struct {
	gorm.Model

	Name string `json:"name"`
}
type TestAC2 struct {
	gorm.Model

	Name string `json:"name"`
	Info string `json:"info"`
}

func (TestAC1) TableName() string { return "tac" }
func (TestAC2) TableName() string { return "tac" }

func TestDialector_AlterColumn(t *testing.T) {
	db, err := gorm.Open(Open("file:testdbac?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	_ = db.AutoMigrate(&TestAC1{})
	_ = db.AutoMigrate(&TestAC2{})
}
