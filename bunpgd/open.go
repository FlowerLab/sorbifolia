package bunpgd

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
	_ "github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"
	"github.com/uptrace/bun/extra/bunslog"
)

func Open(dataSourceName string, opts ...bun.DBOption) (*bun.DB, error) {
	db, err := sql.Open("pg", dataSourceName)
	if err != nil {
		return nil, err
	}
	return bun.NewDB(db, New(), opts...), nil
}

func WithMaxOpenConns(n int) bun.DBOption { return func(db *bun.DB) { db.SetMaxOpenConns(n) } }
func WithMaxIdleConns(n int) bun.DBOption { return func(db *bun.DB) { db.SetMaxIdleConns(n) } }
func WithConnMaxIdleTime(d time.Duration) bun.DBOption {
	return func(db *bun.DB) { db.SetConnMaxIdleTime(d) }
}

func WithOTEL(option ...bunotel.Option) bun.DBOption {
	return func(db *bun.DB) { db.WithQueryHook(bunotel.NewQueryHook(option...)) }
}

func WithSLog(opts ...bunslog.Option) bun.DBOption {
	if len(opts) == 0 {
		opts = []bunslog.Option{
			bunslog.WithQueryLogLevel(slog.LevelDebug),
			bunslog.WithSlowQueryLogLevel(slog.LevelWarn),
			bunslog.WithErrorQueryLogLevel(slog.LevelError),
			bunslog.WithSlowQueryThreshold(3 * time.Second),
		}
	}
	return func(db *bun.DB) { db.WithQueryHook(bunslog.NewQueryHook(opts...)) }
}

func WithCreateTable(ctx context.Context, cancel context.CancelCauseFunc, model ...any) bun.DBOption {
	return func(db *bun.DB) {
		db.RegisterModel(model...)
		for _, v := range model {
			if _, err := db.NewCreateTable().IfNotExists().Model(v).Exec(ctx); err != nil {
				cancel(err)
				break
			}
		}
	}
}
