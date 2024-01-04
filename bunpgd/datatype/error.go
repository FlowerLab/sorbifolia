package datatype

import (
	"errors"
)

var (
	ErrNotSupportValueType           = errors.New("not support value type")
	ErrUnknownTypeCannotMatchScanner = errors.New("unknown type cannot match scanner")
)
