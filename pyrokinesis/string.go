package pyrokinesis

import (
	"unsafe"
)

func (_String) Copy(s string) string {
	return string(String.ToBytes(s))
}

func (_String) ToBytes(s string) []byte {
	//nolint:all
	b := unsafe.StringData(s)
	if b == nil {
		return nil
	}

	//nolint:all
	return unsafe.Slice(b, len(s))
}
