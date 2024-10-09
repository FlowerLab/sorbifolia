package reflectype

import (
	"database/sql"
	"reflect"
)

var (
	NullTime   = reflect.TypeFor[sql.NullTime]()
	NullBool   = reflect.TypeFor[sql.NullBool]()
	NullFloat  = reflect.TypeFor[sql.NullFloat64]()
	NullInt64  = reflect.TypeFor[sql.NullInt64]()
	NullInt32  = reflect.TypeFor[sql.NullInt32]()
	NullInt16  = reflect.TypeFor[sql.NullInt16]()
	NullString = reflect.TypeFor[sql.NullString]()
)
