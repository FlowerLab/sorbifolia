package reflectype

import (
	"database/sql"
	"reflect"
)

var (
	NullTime   = reflect.TypeOf((*sql.NullTime)(nil)).Elem()
	NullBool   = reflect.TypeOf((*sql.NullBool)(nil)).Elem()
	NullFloat  = reflect.TypeOf((*sql.NullFloat64)(nil)).Elem()
	NullInt64  = reflect.TypeOf((*sql.NullInt64)(nil)).Elem()
	NullInt32  = reflect.TypeOf((*sql.NullInt32)(nil)).Elem()
	NullInt16  = reflect.TypeOf((*sql.NullInt16)(nil)).Elem()
	NullString = reflect.TypeOf((*sql.NullString)(nil)).Elem()
)
