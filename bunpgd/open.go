package bunpgd

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

func Open(dataSourceName string, opts ...bun.DBOption) (*bun.DB, error) {
	db, err := sql.Open("pg", dataSourceName)
	if err != nil {
		return nil, err
	}
	return bun.NewDB(db, New(), opts...), nil
}

func WithMaxOpenConns(n int) bun.DBOption { return func(db *bun.DB) { db.DB.SetMaxOpenConns(n) } }
func WithMaxIdleConns(n int) bun.DBOption { return func(db *bun.DB) { db.DB.SetMaxIdleConns(n) } }
func WithConnMaxIdleTime(d time.Duration) bun.DBOption {
	return func(db *bun.DB) { db.DB.SetConnMaxIdleTime(d) }
}
