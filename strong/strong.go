package strong

import (
	"strconv"
)

func Parse[T Type](str string) (T, error) {
	arg := new(T)
	err := _parse(str, arg)
	return *arg, err
}

func Format[T Type](t T) string {
	return _format(t)
}

func Append[T Type](dst []byte, t T) []byte {
	return _append(dst, t)
}

func _parse(str string, data any) error {
	switch data := data.(type) {
	case *bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		*data = val
	case *int:
		val, err := strconv.ParseInt(str, 10, 0)
		if err != nil {
			return err
		}
		*data = int(val)
	case *int8:
		val, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return err
		}
		*data = int8(val)
	case *int16:
		val, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return err
		}
		*data = int16(val)
	case *int32:
		val, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return err
		}
		*data = int32(val)
	case *int64:
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*data = val
	case *uint:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*data = uint(val)
	case *uint8:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*data = uint8(val)
	case *uint16:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*data = uint16(val)
	case *uint32:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*data = uint32(val)
	case *uint64:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*data = val
	default:
		panic("internal error")
	}
	return nil
}

func _format(data any) string {
	switch data := data.(type) {
	case bool:
		if data {
			return "true"
		}
		return "false"
	case int:
		return strconv.FormatInt(int64(data), 10)
	case int8:
		return strconv.FormatInt(int64(data), 10)
	case int16:
		return strconv.FormatInt(int64(data), 10)
	case int32:
		return strconv.FormatInt(int64(data), 10)
	case int64:
		return strconv.FormatInt(data, 10)
	case uint:
		return strconv.FormatUint(uint64(data), 10)
	case uint8:
		return strconv.FormatUint(uint64(data), 10)
	case uint16:
		return strconv.FormatUint(uint64(data), 10)
	case uint32:
		return strconv.FormatUint(uint64(data), 10)
	case uint64:
		return strconv.FormatUint(data, 10)
	}
	panic("internal error")
}

func _append(dst []byte, data any) []byte {
	switch data := data.(type) {
	case bool:
		if data {
			return append(dst, "true"...)
		}
		return append(dst, "false"...)
	case int:
		return strconv.AppendInt(dst, int64(data), 10)
	case int8:
		return strconv.AppendInt(dst, int64(data), 10)
	case int16:
		return strconv.AppendInt(dst, int64(data), 10)
	case int32:
		return strconv.AppendInt(dst, int64(data), 10)
	case int64:
		return strconv.AppendInt(dst, data, 10)
	case uint:
		return strconv.AppendUint(dst, uint64(data), 10)
	case uint8:
		return strconv.AppendUint(dst, uint64(data), 10)
	case uint16:
		return strconv.AppendUint(dst, uint64(data), 10)
	case uint32:
		return strconv.AppendUint(dst, uint64(data), 10)
	case uint64:
		return strconv.AppendUint(dst, data, 10)
	}
	panic("internal error")
}
