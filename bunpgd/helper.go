package bunpgd

import (
	"encoding/json"
	"reflect"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/datatype"
)

type Array struct {
	rt reflect.Type
	rv reflect.Value
}

func (a *Array) AppendQuery(gen schema.QueryGen, b []byte) ([]byte, error) {
	return datatype.TypeAppender(a.rt)(gen, b, a.rv), nil
}

func ToArray[T any](arr []T) *Array {
	rv := reflect.ValueOf(arr)
	return ArrayFormReflect(rv.Type(), rv)
}

func ArrayFormReflect(rt reflect.Type, rv reflect.Value) *Array {
	return &Array{rv: rv, rt: rt}
}

type JSON struct {
	rv reflect.Value
}

func (j *JSON) AppendQuery(gen schema.QueryGen, b []byte) ([]byte, error) {
	bb, err := json.Marshal(j.rv.Interface())
	if err != nil {
		return dialect.AppendError(b, err), nil
	}

	return gen.Dialect().AppendJSON(b, bb), nil
}

func JSONFormReflect(rv reflect.Value) *JSON {
	return &JSON{rv: rv}
}
