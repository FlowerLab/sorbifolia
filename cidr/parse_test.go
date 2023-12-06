package cidr

import (
	"testing"
)

var testParseGroup = []struct {
	str []string
	err bool
}{
	{[]string{}, false},
	{[]string{"1.1.1.1"}, false},
	{[]string{"1.1.1.0/24"}, false},
	{[]string{"1.1.1.1-1.1.1.10"}, false},
	{[]string{"1.3.1.1", "1.2.1.0/24", "1.1.1.1-1.1.1.10"}, false},
	{[]string{"1.3.1.1 ", " 1.2.1.0/24", "1.1.1.1 - 1.1.1.10"}, false},

	{[]string{""}, true},
	{[]string{"1.1.1.257"}, true},
	{[]string{"1.1.1.0/33"}, true},
	{[]string{"1.1.1.1-1.1.1.257"}, true},
	{[]string{"1.1.1.x-1.1.1.23"}, true},
}

func TestParseGroup(t *testing.T) {
	for _, v := range testParseGroup {
		if _, err := ParseGroup(v.str); err != nil && !v.err {
			t.Errorf("parse %s err: %s", v.str, err)
		}
	}
}
