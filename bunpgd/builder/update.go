package builder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uptrace/bun"
)

// OptionalUpdate updates the fields of the target struct with values from the update struct.
//
// If field in the struct is a pointer and is nil, it is not updated, else is always updated.
func OptionalUpdate(q *bun.UpdateQuery, v any, skip ...string) *bun.UpdateQuery {
	needSkip := func(key string) bool {
		for _, s := range skip {
			if key == s {
				return true
			}
		}
		return false
	}

	var (
		rv = reflect.Indirect(reflect.ValueOf(v))
		rt = rv.Type()
	)
	if rt.Kind() != reflect.Struct {
		return q.Err(fmt.Errorf("expected a struct, got %T", v))
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Anonymous || !field.IsExported() {
			continue
		}

		tag, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if tag == "-" || needSkip(tag) {
			continue
		}

		var (
			kind = field.Type.Kind()
			val  = rv.Field(i)
		)

		switch kind {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128,
			reflect.Array, reflect.Map, reflect.Slice, reflect.String, reflect.Struct:
			q.Set("? = ?", bun.Ident(tag), val.Interface())

		case reflect.Pointer:
			if val.IsNil() {
				continue
			}

			if field.Type.Elem().Kind() != reflect.Struct {
				q.Set("? = ?", bun.Ident(tag), val.Elem().Interface())
			} else {
				q.Set("? = ?", bun.Ident(tag), val.Interface())
			}

		default:
			return q.Err(fmt.Errorf("unexpected data type %s", kind))
		}
	}

	return q
}
