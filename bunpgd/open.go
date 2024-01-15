package bunpgd

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
)

func Open(dataSourceName string, opts ...bun.DBOption) (*bun.DB, error) {
	db, err := sql.Open("pgx/v5", dataSourceName)
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
