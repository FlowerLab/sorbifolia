package datatype

import (
	"database/sql"
	"database/sql/driver"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type gormItr interface {
	driver.Valuer
	sql.Scanner
	GormDataType() string
	GormDBDataType(*gorm.DB, *schema.Field) string
}

func isPostgres(db *gorm.DB, res string) string {
	switch db.Name() {
	case "postgres":
		return res
	}
	return ""
}
