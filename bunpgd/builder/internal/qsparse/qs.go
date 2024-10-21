package qsparse

import (
	"errors"
	"net/url"
	"reflect"

	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/query"
)

func QS(qs url.Values, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return errors.New("value must be a pointer")
	}
	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	rv = rv.Elem()

	var (
		rt = rv.Type()
		sf = query.CachedTypeFields(rt)
	)

	for key, val := range qs {
		field := sf.List[sf.NameIndex[key]]
		if field.Name != key {
			continue
		}
		if len(val) == 0 {
			continue
		}

		if err := from(val, field, rv.FieldByIndex(field.Index)); err != nil {
			return err
		}
	}

	return nil
}
