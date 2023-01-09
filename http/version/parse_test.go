package version

import (
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

	for i := 0; i < len(testVersion); i++ {
		version := testVersion[i]
		major, minor, ok := parseHTTPVersion(version.Bytes())
		if !ok {
			t.Error("parse error")
		}

		if !reflect.DeepEqual(version.Major, major) {
			t.Errorf("expected %v,got %v", version.Major, major)
		}
		if !reflect.DeepEqual(version.Minor, minor) {
			t.Errorf("expected %v,got %v", version.Minor, minor)
		}
	}
}
