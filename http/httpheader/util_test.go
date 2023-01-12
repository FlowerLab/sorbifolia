package httpheader

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCleanSuffixSpace(t *testing.T) {
	test := [][]byte{
		[]byte("ddtest123"),
		[]byte("asd    123    "),
		[]byte("        "),
		[]byte("dsada123    "),
	}
	ans := [][]byte{
		[]byte("ddtest123"),
		[]byte("asd    123"),
		nil,
		[]byte("dsada123"),
	}

	for i, v := range test {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			nv := cleanSuffixSpace(v)
			if !reflect.DeepEqual(ans[i], nv) {
				t.Errorf("expected %v,got %v", string(ans[i]), string(nv))
			}
		})
	}
}
