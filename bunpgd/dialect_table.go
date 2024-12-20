package bunpgd

import (
	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/datatype"
	"go.x2ox.com/sorbifolia/bunpgd/sqltype"
)

func (d *Dialect) OnTable(table *schema.Table) {
	for _, field := range table.FieldMap {
		sqltype.Set(field)
		datatype.Set(field)
	}
}
