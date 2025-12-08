package builder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uptrace/bun"
	"go.x2ox.com/sorbifolia/bunpgd"
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

func OptionalForceUpdate(q *bun.UpdateQuery, v any, force, skip []string) *bun.UpdateQuery {
	needSkip := func(key string) bool {
		for _, s := range skip {
			if key == s {
				return true
			}
		}
		return false
	}

	isForce := func(key string) bool {
		for _, s := range force {
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
				if isForce(tag) {
					q.Set("? = NULL", bun.Ident(tag))
				}
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

func SelectUpdate(q *bun.UpdateQuery, v any, selectKey ...string) *bun.UpdateQuery {
	if len(selectKey) == 0 {
		return q
	}

	has := func(key string) bool {
		for _, s := range selectKey {
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
		if tag == "-" || !has(tag) {
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
				q.Set("? = NULL", bun.Ident(tag))
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

type Updater struct {
	q *bun.UpdateQuery
	v any

	pbc  bool  // protobuf compatible: Updater.key is protobuf tag json key, bun.UpdateQuery.Set use json key
	mode uint8 // 0: ignore mode, 1: select mode
	key  []string
}

func UseUpdater(q *bun.UpdateQuery, v any) *Updater { return &Updater{q: q, v: v} }
func (x *Updater) Ignore(arr []string) *Updater     { x.key, x.mode = arr, 0; return x }
func (x *Updater) Select(arr []string) *Updater     { x.key, x.mode = arr, 1; return x }
func (x *Updater) PB() *Updater                     { x.pbc = true; return x }
func (x *Updater) Exec() *bun.UpdateQuery           { return x.q }

func (x *Updater) parseKey(tag reflect.StructTag) (key string, sqlKey bun.Ident) {
	if key, _, _ = strings.Cut(tag.Get("json"), ","); key == "-" {
		return "", "" // json:"search_key,omitempty"
	}

	sqlKey = bun.Ident(key)

	if !x.pbc {
		return
	}

	ps := tag.Get("protobuf")
	if ps == "" {
		return // protobuf:"bytes,2,opt,name=search_key,json=searchKey,proto3,oneof"
	}

	si := strings.Index(ps, "json=") + 5 // len("json=") == 5
	if si == 4 {
		return
	}

	if ei := strings.Index(ps[si:], ","); ei == -1 {
		key = ps[si:]
	} else {
		key = ps[si : si+ei]
	}

	return
}

func (x *Updater) has(key string) bool {
	for _, k := range x.key {
		if k == key {
			return true
		}
	}

	return false
}

func (x *Updater) check(key string) bool {
	if key == "" {
		return false
	}

	switch x.mode {
	case 0: // ignore mode
		return !x.has(key)
	case 1: // select mode
		return x.has(key)
	default:
		return false
	}
}

func (x *Updater) exec() *bun.UpdateQuery {
	var (
		q  = x.q
		rv = reflect.Indirect(reflect.ValueOf(x.v))
		rt = rv.Type()
	)

	if rt.Kind() != reflect.Struct {
		return q.Err(fmt.Errorf("expected a struct, got %T", x.v))
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Anonymous || !field.IsExported() {
			continue
		}

		key, sqlKey := x.parseKey(field.Tag)
		if !x.check(key) {
			continue
		}

		var (
			kind = field.Type.Kind()
			val  = rv.Field(i)
		)

		switch kind {
		case reflect.Slice:
			if val.IsNil() {
				q.Set("? = NULL", sqlKey)
				continue
			}
			q.Set("? = ?", sqlKey, bunpgd.ArrayFormReflect(field.Type, val))

		case reflect.Map:
			if val.IsNil() {
				q.Set("? = NULL", sqlKey)
				continue
			}
			q.Set("? = ?", sqlKey, val.Interface())

		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128,
			reflect.Array, reflect.String, reflect.Struct:
			q.Set("? = ?", sqlKey, val.Interface())

		case reflect.Pointer:
			switch {
			case val.IsNil():
				q.Set("? = NULL", sqlKey)
			case field.Type.Elem().Kind() == reflect.Struct:
				q.Set("? = ?", sqlKey, val.Interface())
			default:
				q.Set("? = ?", sqlKey, val.Elem().Interface())
			}

		default:
			return q.Err(fmt.Errorf("unexpected data type %s", kind))
		}
	}

	return q
}
