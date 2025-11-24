package bunpgd

import (
	"reflect"

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
