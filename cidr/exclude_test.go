package cidr

import (
	"math/big"
	"net/netip"
	"strings"
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
		contains: must(ParseRange, "1.0.0.0-1.0.1.100"),
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

var testExcludeLength = []struct {
	include CIDR
	exclude Group
	expect  *big.Int
}{
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		expect:  big.NewInt(155),
	},
}

func TestExclude_Length(t *testing.T) {
	for _, val := range testExcludeLength {
		e := &Exclude{e: val.exclude, i: val.include}
		b := e.Length()
		if b.Cmp(val.expect) != 0 {
			t.Errorf("expected value is %s, but got %s", val.expect.String(), b.String())
		}
	}
}

var testExcludeNextIP = []struct {
	include CIDR
	exclude Group
	ip      netip.Addr
	next    netip.Addr
}{
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100")}},
		ip:      netip.IPv4Unspecified(),
		next:    netip.AddrFrom4([4]byte{1, 0, 0, 0}),
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 1}),
		next:    netip.AddrFrom4([4]byte{1, 0, 0, 2}),
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 3}),
		next:    netip.AddrFrom4([4]byte{1, 0, 0, 101}),
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 255}),
		next:    netip.IPv4Unspecified(),
	},
}

func TestExclude_NextIP(t *testing.T) {
	for _, val := range testExcludeNextIP {
		e := &Exclude{e: val.exclude, i: val.include}
		if b := e.NextIP(val.ip); b.Compare(val.next) != 0 {
			t.Errorf("expected value is %s, but got %s", val.next.String(), b.String())
		}
	}
}

var testExcludeString = []struct {
	exclude Group
	dst     []string
}{
	{
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		dst:     []string{"1.0.0.0-1.0.0.100"},
	},
}

func TestExclude_String(t *testing.T) {
	for _, val := range testExcludeString {
		e := &Exclude{e: val.exclude}
		dst := e.Strings()

		if strings.Join(dst, "|") != strings.Join(val.dst, "|") {
			t.Errorf("expected value is %s, but got %s", val.dst, dst)
		}
	}
}