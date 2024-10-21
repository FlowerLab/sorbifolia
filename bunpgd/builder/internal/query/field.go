package query

import (
	"reflect"

	"github.com/uptrace/bun"
	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/attribute"
	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/flag"
	op "go.x2ox.com/sorbifolia/bunpgd/builder/internal/operator"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

// A Field represents a single Field found in a struct.
type Field struct {
	Name  string
	Index []int
	Typ   reflect.Type

	Flag flag.Flag
	Op   op.Operator
	Attr attribute.Attribute
	Key  bun.Ident
}

func (f Field) ReflectQuery(v reflect.Value) reflectype.BunQueryBuilder {
	if f.Flag.And(reflectype.QueryBuilder) {
		return HandleFunc(func(q bun.QueryBuilder) bun.QueryBuilder {
			v.MethodByName("BunQueryBuilder").Call([]reflect.Value{reflect.ValueOf(q)})
			return q
		})
	}
	if v.IsZero() || (f.Flag.And(reflect.Pointer) && v.Elem().IsZero()) {
		return nil
	}

	return f.handle(v)
}

func (f Field) handle(v reflect.Value) HandleFunc {
	switch f.Op {
	case op.In, op.NotIn:
		v = reflect.Indirect(v)

		switch v.Len() {
		case 0:
			return nil
		case 1: // use equal
		default:
			return where("? ? (?)", f.Key, f.Op, bun.In(v.Interface()))
		}

	case op.NotLike, op.Like:
		return f.handleLike(reflect.Indirect(v))

	case op.IsDistinct, op.IsNotDistinct, op.IsNull, op.IsNotNull, op.IsTrue, op.IsNotTrue,
		op.IsFalse, op.IsNotFalse, op.IsUnknown, op.IsNotUnknown:
		return where("? ?", f.Key, f.Op)
	}

	return where("? ? ?", f.Key, f.Op, v.Interface())
}

func (f Field) handleLike(v reflect.Value) HandleFunc {
	if f.Flag.And(reflect.Slice) {
		length := v.Len()
		switch length {
		case 0:
			return nil
		case 1:
			v = v.Index(0)
		default:
			return func(q bun.QueryBuilder) bun.QueryBuilder {
				return q.WhereGroup(" AND ", func(q bun.QueryBuilder) bun.QueryBuilder {
					for i := range length {
						if str := v.Index(i); !str.IsZero() {
							q.WhereOr("? ? ?", f.Key, f.Op, f.Attr.FormatValue(f.Op, str))
						}
					}
					return q
				})
			}
		}
	}

	return where("? ? ?", f.Key, f.Op, f.Attr.FormatValue(f.Op, v))
}

func where(qs string, args ...any) HandleFunc {
	return func(q bun.QueryBuilder) bun.QueryBuilder { return q.Where(qs, args...) }
}
