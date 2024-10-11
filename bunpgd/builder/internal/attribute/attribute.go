package attribute

import (
	"reflect"

	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/operator"
)

type Attribute string

func (x Attribute) FormatValue(op operator.Operator, v reflect.Value) any {
	switch op {
	case operator.Like, operator.NotLike:
		return x.formatLikeValue(v)
	default:
		panic("not support operator")
	}
}

func (x Attribute) formatLikeValue(v reflect.Value) string {
	str := v.Interface().(string)

	switch x {
	case "left", "l", "L":
		return "%" + str
	case "right", "r", "R":
		return str + "%"
	default:
		return "%" + str + "%"
	}
}
