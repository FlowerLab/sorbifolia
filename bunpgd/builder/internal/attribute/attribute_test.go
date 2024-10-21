package attribute

import (
	"reflect"
	"testing"

	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/operator"
)

var testAttrData = []struct {
	attr  Attribute
	op    operator.Operator
	val   any
	res   any
	panic bool
}{
	{attr: "left", op: operator.NotLike, val: "x", res: "%x"},
	{attr: "L", op: operator.NotLike, val: "x", res: "%x"},
	{attr: "right", op: operator.NotLike, val: "x", res: "x%"},
	{attr: "R", op: operator.NotLike, val: "x", res: "x%"},
	{attr: "", op: operator.NotLike, val: "x", res: "%x%"},

	{op: operator.In, panic: true},
}

func TestAttribute_FormatValue(t *testing.T) {
	for _, data := range testAttrData {
		t.Run(data.op.String(), func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					if !data.panic {
						t.Errorf("panic expected: %t, got: %t", data.panic, err)
					}
				}
			}()

			val := data.attr.FormatValue(data.op, reflect.ValueOf(data.val))
			if !reflect.DeepEqual(val, data.res) {
				t.Errorf("FormatValue(%v, %v): got %v, want %v", data.op, val, data.res, val)
			}
		})
	}
}
