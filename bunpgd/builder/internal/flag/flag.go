package flag

import (
	"reflect"

	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

type Flag uint16

var (
	flag = [13]any{
		reflect.Pointer, reflect.Slice,

		reflectype.QueryBuilder, reflectype.FromQS,
		reflectype.JSONMarshaler, reflectype.JSONUnmarshaler,

		reflectype.Time, reflectype.TimeDuration,
		reflectype.IP, reflectype.IPNet,
		reflectype.HardwareAddr,
		reflectype.Addr, reflectype.Prefix,
	}

	bitmap = map[any]Flag{
		flag[0]: 1 << 0, flag[1]: 1 << 1,

		flag[2]: 1 << 2, flag[3]: 1 << 3,
		flag[4]: 1 << 4, flag[5]: 1 << 5,

		flag[6]: 1 << 6, flag[7]: 1 << 7,
		flag[8]: 1 << 8, flag[9]: 1 << 9,
		flag[10]: 1 << 10,
		flag[11]: 1 << 11, flag[12]: 1 << 12,
	}
)

func Bit(a any) Flag           { return bitmap[a] }
func (f *Flag) Has(a any) bool { return *f&Bit(a) != 0 }

func From(rt reflect.Type) Flag {
	var f Flag

	for _, v := range flag[2:6] {
		if rt.Implements(v.(reflect.Type)) {
			f |= bitmap[v]
		}
	}

	kind := rt.Kind()
	if kind == reflect.Pointer {
		f |= 1
		rt = rt.Elem()
		kind = rt.Kind()
	}
	if kind == reflect.Slice {
		f |= 2
	}

	for _, v := range flag[6:] {
		if rt == v {
			f |= bitmap[v]
		}
	}

	return f
}
