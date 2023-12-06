package cidr

import (
	"testing"
)

func must[T any](fn func(s string) (*T, error), s string) *T {
	t, err := fn(s)
	if err != nil {
		panic(err)
	}
	return t
}

var testExcludeContains = []struct {
	include  CIDR
	exclude  Group
	contains CIDR
	status   ContainsStatus
}{
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		contains: must(ParsePrefix, "1.0.0.0/25"),
		status:   ContainsPartially,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		contains: must(ParsePrefix, "1.0.0.0/26"),
		status:   ContainsNot,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		contains: must(ParsePrefix, "1.0.0.233/32"),
		status:   Contains,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		contains: &Group{}, // unhandled behavior
		status:   ContainsNot,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		contains: &Group{arr: []Consecutive{must(ParseRange, "1.0.0.111-1.0.0.112")}},
		status:   ContainsNot, // unhandled behavior
	},

	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParsePrefix, "1.0.0.0/30")}},
		contains: must(ParsePrefix, "1.0.0.0/25"),
		status:   ContainsPartially,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParsePrefix, "1.0.0.0/30")}},
		contains: must(ParsePrefix, "1.0.0.0/31"),
		status:   ContainsNot,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParsePrefix, "1.0.0.0/30")}},
		contains: must(ParsePrefix, "1.0.0.233/32"),
		status:   Contains,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParsePrefix, "1.0.0.0/30")}},
		contains: &Group{}, // unhandled behavior
		status:   ContainsNot,
	},

	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseSingle, "1.0.0.1")}},
		contains: must(ParsePrefix, "1.0.0.0/25"),
		status:   ContainsPartially,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseSingle, "1.0.0.1")}},
		contains: must(ParsePrefix, "1.0.0.6/31"),
		status:   Contains,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseSingle, "1.0.0.1")}},
		contains: must(ParsePrefix, "1.0.0.233/32"),
		status:   Contains,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseSingle, "1.0.0.1")}},
		contains: must(ParseSingle, "1.0.0.1"),
		status:   ContainsNot,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseSingle, "1.0.0.1")}},
		contains: must(ParseSingle, "1.0.0.13"),
		status:   Contains,
	},
}

func TestExclude_Contains(t *testing.T) {
	for _, val := range testExcludeContains {
		e := &Exclude{e: val.exclude, i: val.include}
		if status := e.Contains(val.contains); status != val.status {
			t.Errorf("expected value is %d, but got %d", val.status, status)
		}
	}
}
