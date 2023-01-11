package version

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		v          Version
		is09, null bool
	}{
		{v: Version{Major: 0, Minor: 0}, is09: false, null: true},
		{v: Version{Major: 0, Minor: 9}, is09: true, null: false},
		{v: Version{Major: 1, Minor: 0}, is09: false, null: false},
		{v: Version{Major: 1, Minor: 1}, is09: false, null: false},
		{v: Version{Major: 2, Minor: 0}, is09: false, null: false},
		{v: Version{Major: 3, Minor: 0}, is09: false, null: false},
		{v: Version{Major: 4, Minor: 0}, is09: false, null: false},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ver, ok := Parse(v.v.Bytes())
			if !ok {
				t.Error("parse error")
			}

			if !reflect.DeepEqual(v.v, ver) {
				t.Errorf("expected %v, got %v\n", v.v, ver)
			}
			if v.null != ver.Null() {
				t.Errorf("expected %v, got %v\n", v.null, ver.Null())
			}
			if v.is09 != ver.Is09() {
				t.Errorf("expected %v, got %v\n", v.is09, ver.Is09())
			}
		})
	}
}
