package datatype

import (
	"database/sql/driver"
	"fmt"

	"go.x2ox.com/sorbifolia/strong"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Array[T strong.Type] []T

// GormDataType gorm common data type
func (Array[T]) GormDataType() string { return "Array" }

// GormDBDataType gorm db data type
func (Array[T]) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "text[]"
	}
	return ""
}

// Scan implements the sql.Scanner interface.
func (a *Array[T]) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("cannot convert %T to Array", src)
}

func (a *Array[T]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, a.GormDataType())
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make([]T, len(elems))
		for i, v := range elems {
			if b[i], err = strong.Parse[T](string(v)); err != nil {
				return fmt.Errorf("parsing array element index %d: %v", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Array[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		b := make([]byte, 1, 1+2*n) // {} and N + N-1 bytes of delimiters
		b[0] = '{'

		b = strong.Append(b, a[0])
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strong.Append(b, a[i])
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}
