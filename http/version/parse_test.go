package version

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	testVersion := []Version{
		{Major: 0, Minor: 9},
		{Major: 1, Minor: 0},
		{Major: 1, Minor: 1},
		{Major: 2, Minor: 0},
		{Major: 3, Minor: 0},
	}

	for i, v := range testVersion {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			version, ok := Parse(v.Bytes())
			if !ok {
				t.Error("parse error")
			}

			if !reflect.DeepEqual(v.Major, version.Major) {
				t.Errorf("expected %v,got %v", v.Major, version.Major)
			}
			if !reflect.DeepEqual(v.Minor, version.Minor) {
				t.Errorf("expected %v,got %v", v.Minor, version.Minor)
			}
		})
	}
}
