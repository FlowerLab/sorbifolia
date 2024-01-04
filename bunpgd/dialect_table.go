package bunpgd

import (
	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/datatype"
	"go.x2ox.com/sorbifolia/bunpgd/sqltype"
)

func (d *Dialect) OnTable(table *schema.Table) {
	for _, field := range table.FieldMap {
		d.onField(field)
	}
}

func (d *Dialect) onField(field *schema.Field) {
	if err := sqltype.Set(field); err != nil {
		panic(err)
	}
	datatype.Set(field)
}
