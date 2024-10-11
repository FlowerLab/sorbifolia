package builder

import (
	"reflect"

	"github.com/uptrace/bun"
)

func Query(q *bun.SelectQuery, v any) {
	var (
		rv = reflect.Indirect(reflect.ValueOf(v))
		rt = rv.Type()
	)

	tf := CachedTypeFields(rt)
	for _, field := range tf.List {
		fv := rv.FieldByIndex(field.Index)
		itr := field.ReflectQuery(fv)

		if itr == nil {
			continue
		}

		q.ApplyQueryBuilder(itr.BunQueryBuilder)
	}
}

type HandleFunc func(bun.QueryBuilder) bun.QueryBuilder

func (f HandleFunc) BunQueryBuilder(q bun.QueryBuilder) bun.QueryBuilder { return f(q) }
