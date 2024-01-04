package datatype

import (
	"fmt"
	"reflect"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/schema"
)

type Adapter struct {
	SQLType string
	Type    reflect.Type

	Append schema.AppenderFunc
	Scan   schema.ScannerFunc
	IsZero schema.IsZeroerFunc
}

func (a *Adapter) Set(field *schema.Field) {
	if a.Append != nil {
		field.Append = a.Append
	}
	if a.Scan != nil {
		field.Scan = a.Scan
	}
	if a.IsZero != nil {
		field.IsZero = a.IsZero
	}
}

func (a *Adapter) Ptr() *Adapter {
	b := &Adapter{IsZero: a.IsZero}
	if a.Append != nil {
		b.Append = schema.PtrAppender(a.Append)
	}
	if a.Scan != nil {
		b.Scan = schema.PtrScanner(a.Scan)
	}
	return b
}

func (a *Adapter) Addr() *Adapter {
	b := &Adapter{IsZero: a.IsZero}
	if a.Append != nil {
		b.Append = addrAppender(a.Append)
	}
	if a.Scan != nil {
		b.Scan = addrScanner(a.Scan)
	}
	return b
}

func addrScanner(fn schema.ScannerFunc) schema.ScannerFunc {
	return func(dest reflect.Value, src any) error {
		if !dest.CanAddr() {
			return fmt.Errorf("bunpgd: Scan(nonaddressable %T)", dest.Interface())
		}
		if err := fn(dest.Addr(), src); err != nil {
			return err
		}

		if dest.Elem().IsZero() {
			dest.SetZero()
		}
		return nil
	}
}

func addrAppender(fn schema.AppenderFunc) schema.AppenderFunc {
	return func(fmter schema.Formatter, b []byte, v reflect.Value) []byte {
		if !v.CanAddr() {
			err := fmt.Errorf("bunpgd: Append(nonaddressable %T)", v.Interface())
			return dialect.AppendError(b, err)
		}
		return fn(fmter, b, v.Addr())
	}
}
