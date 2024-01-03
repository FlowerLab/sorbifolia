package b2s

import (
	"unsafe"
)

func S(s string) []byte { return unsafe.Slice(unsafe.StringData(s), len(s)) }
func B(b []byte) string { return unsafe.String(unsafe.SliceData(b), len(b)) }

func A(a any) []byte {
	switch a := a.(type) {
	case nil:
		return nil
	case []byte:
		return a
	case string:
		return S(a)
	default:
		panic("not support datatype")
	}
}
