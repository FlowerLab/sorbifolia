package bunpgd

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

func init() {
	if bun.Version() != "1.1.16" {
		panic("not support version: " + bun.Version())
	}
}

var (
	_ schema.QueryAppender = nil
	_ schema.AppenderFunc  = nil
	_ schema.ScannerFunc   = nil
	_ sql.Scanner          = nil
)
