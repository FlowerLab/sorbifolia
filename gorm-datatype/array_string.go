package datatype

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ArrayString []string

// GormDataType gorm common data type
func (ArrayString) GormDataType() string { return "ArrayString" }

// GormDBDataType gorm db data type
func (ArrayString) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "text[]"
	}
	return ""
}

// Scan implements the sql.Scanner interface.
func (a *ArrayString) Scan(src any) error {
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

func (a *ArrayString) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, a.GormDataType())
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make([]string, len(elems))
		for i, v := range elems {
			if b[i] = string(v); v == nil {
				return fmt.Errorf("parsing array element index %d: cannot convert nil to string", i)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a ArrayString) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		b := make([]byte, 1, 1+3*n) // {} and 2*N + N-1 bytes of delimiters
		b[0] = '{'

		b = appendArrayQuotedBytes(b, []byte(a[0]))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, []byte(a[i]))
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}
