package bunpgd

import (
	"database/sql"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/feature"
	"github.com/uptrace/bun/schema"
)

type Dialect struct {
	schema.BaseDialect

	tables *schema.Tables
}

func New() *Dialect {
	d := new(Dialect)
	d.tables = schema.NewTables(d)
	return d
}

func (d *Dialect) Init(_ *sql.DB)         {}
func (d *Dialect) Name() dialect.Name     { return dialect.PG }
func (d *Dialect) DefaultSchema() string  { return "public" }
func (d *Dialect) IdentQuote() byte       { return '"' }
func (d *Dialect) DefaultVarcharLen() int { return 0 }
func (d *Dialect) Tables() *schema.Tables { return d.tables }
func (d *Dialect) AppendSequence(b []byte, _ *schema.Table, _ *schema.Field) []byte {
	return append(b, " GENERATED BY DEFAULT AS IDENTITY"...)
}

func (d *Dialect) Features() feature.Feature {
	return feature.CTE |
		feature.WithValues |
		feature.Returning |
		feature.InsertReturning |
		feature.DefaultPlaceholder |
		feature.DoubleColonCast |
		feature.InsertTableAlias |
		feature.UpdateTableAlias |
		feature.DeleteTableAlias |
		feature.TableCascade |
		feature.TableIdentity |
		feature.TableTruncate |
		feature.TableNotExists |
		feature.InsertOnConflict |
		feature.SelectExists |
		feature.GeneratedIdentity |
		feature.CompositeIn |
		feature.DeleteReturning
}
