package builder

import (
	"iter"
	"reflect"

	"github.com/uptrace/bun"
)

type HandleFunc func(bun.QueryBuilder) bun.QueryBuilder

func (f HandleFunc) BunQueryBuilder(q bun.QueryBuilder) bun.QueryBuilder { return f(q) }

func Generate(v any) iter.Seq[HandleFunc] {
	var (
		rv = reflect.Indirect(reflect.ValueOf(v))
		tf = CachedTypeFields(rv.Type())
	)

	return func(yield func(HandleFunc) bool) {
		for _, field := range tf.List {
			var (
				fv  = rv.FieldByIndex(field.Index)
				itr = field.ReflectQuery(fv)
			)

			if itr == nil {
				continue
			}
			if !yield(itr.BunQueryBuilder) {
				break
			}
		}
	}
}
