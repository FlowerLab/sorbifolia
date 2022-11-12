package pyrokinesis

import (
	"reflect"
)

func Call[T any](s T, name string, args []reflect.Value) []reflect.Value {
	var (
		v = reflect.ValueOf(s)
		p = reflect.Value{}
	)

	if v.Kind() == reflect.Ptr {
		p = v
		v = p.Elem()
	} else {
		p = reflect.ValueOf(new(T))
	}

	if v.Kind() != reflect.Struct {
		panic("T is not struct")
	}

	for i, t := 0, v.Type(); i < t.NumMethod(); i++ {
		if method := t.Method(i); method.Name == name {
			return method.Func.Call(append([]reflect.Value{v}, args...))
		}
	}

	for i, t := 0, p.Type(); i < t.NumMethod(); i++ {
		if method := t.Method(i); method.Name == name {
			return method.Func.Call(append([]reflect.Value{p}, args...))
		}
	}

	panic("method not found")
}
