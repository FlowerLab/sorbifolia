package cidr

import (
	"errors"
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

	{
		include:  unknownContainsStatus{},
		exclude:  Group{arr: []Consecutive{must(ParseSingle, "1.0.0.1")}},
		contains: must(ParseSingle, "1.0.0.13"),
		status:   ContainsNot,
	},
	{
		include:  must(ParseRange, "1.0.0.1-1.0.0.10"),
		exclude:  Group{},
		contains: &Group{arr: []Consecutive{must(ParseSingle, "1.0.0.1")}},
		status:   ContainsNot,
	},
}

type unknownContainsStatus struct {
}

func (unknownContainsStatus) ContainsIP(_ netip.Addr) bool   { panic("implement me") }
func (unknownContainsStatus) NextIP(_ netip.Addr) netip.Addr { panic("implement me") }
func (unknownContainsStatus) Length() *big.Int               { panic("implement me") }
func (unknownContainsStatus) Contains(_ CIDR) ContainsStatus { return 255 }

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
		ip:      invalidIP,
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
		next:    invalidIP,
	},

	{
		include: &Group{},
		exclude: Group{},
		ip:      invalidIP,
		next:    invalidIP,
	},
	{
		include: &Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100")}},
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.80")}},
		ip:      invalidIP,
		next:    netip.AddrFrom4([4]byte{1, 0, 0, 81}),
	},
	{
		include: &Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100")}},
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.80")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 81}),
		next:    netip.AddrFrom4([4]byte{1, 0, 0, 82}),
	},
	{
		include: &Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100")}},
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.80")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 100}),
		next:    invalidIP,
	},
	{
		include: &Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100"), must(ParseRange, "1.0.0.130-1.0.0.140")}},
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.80")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 100}),
		next:    netip.AddrFrom4([4]byte{1, 0, 0, 130}),
	},
	{
		include: &Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.100"), must(ParseRange, "1.0.0.130-1.0.0.140")}},
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.3-1.0.0.80")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 133}),
		next:    netip.AddrFrom4([4]byte{1, 0, 0, 134}),
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

var testExcludeContainsIP = []struct {
	include  CIDR
	exclude  Group
	ip       netip.Addr
	contains bool
}{
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		ip:       netip.AddrFrom4([4]byte{1, 0, 0, 1}),
		contains: false,
	},
	{
		include:  must(ParsePrefix, "1.0.0.0/24"),
		exclude:  Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		ip:       netip.AddrFrom4([4]byte{1, 0, 0, 101}),
		contains: true,
	},
	{
		include:  must(ParseSingle, "1.0.0.2"),
		ip:       netip.AddrFrom4([4]byte{1, 0, 0, 2}),
		contains: true,
	},
}

func TestExclude_ContainsIP(t *testing.T) {
	for _, val := range testExcludeContainsIP {
		e := &Exclude{e: val.exclude, i: val.include}
		if contains := e.ContainsIP(val.ip); contains != val.contains {
			t.Errorf("expected value is %t, but got %t", val.contains, contains)
		}
	}
}

var testExcludeAddAddress = []struct {
	include CIDR
	exclude Group
	ip      netip.Addr
	error   error
}{
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{}},
		ip:      netip.AddrFrom4([4]byte{1, 2, 0, 1}),
		error:   ErrNotInAddressRange,
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 23}),
		error:   ErrHasBeenExcluded,
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.0-1.0.0.100")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 123}),
		error:   nil,
	},
}

func TestExclude_AddAddress(t *testing.T) {
	for _, val := range testExcludeAddAddress {
		e := &Exclude{e: val.exclude, i: val.include}
		err := e.AddAddress(val.ip)
		switch {
		case err == nil && val.error == nil:
			continue

		case err != nil && val.error != nil:
			if !errors.Is(err, val.error) {
				t.Errorf("expected value is %s, but got %s", val.error, err)
			}
		default:
			t.Errorf("expected value is %s, but got %s", val.error, err)
		}
	}
}

var testExcludeDelAddress = []struct {
	include CIDR
	exclude Group
	ip      netip.Addr
	error   error
	dst     []string
}{
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{}},
		ip:      netip.AddrFrom4([4]byte{1, 2, 0, 1}),
		error:   ErrNotInAddressRange,
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.30-1.0.0.40"), must(ParseRange, "1.0.0.50-1.0.0.60")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 90}),
		dst:     []string{"1.0.0.30-1.0.0.40", "1.0.0.50-1.0.0.60"},
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.30-1.0.0.30"), must(ParseRange, "1.0.0.50-1.0.0.60")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 30}),
		dst:     []string{"1.0.0.50-1.0.0.60"},
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.30-1.0.0.40"), must(ParseRange, "1.0.0.50-1.0.0.60")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 30}),
		dst:     []string{"1.0.0.31-1.0.0.40", "1.0.0.50-1.0.0.60"},
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.30-1.0.0.40"), must(ParseRange, "1.0.0.50-1.0.0.60")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 40}),
		dst:     []string{"1.0.0.30-1.0.0.39", "1.0.0.50-1.0.0.60"},
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseRange, "1.0.0.30-1.0.0.40")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 33}),
		dst:     []string{"1.0.0.30-1.0.0.32", "1.0.0.34-1.0.0.40"},
	},

	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseSingle, "1.0.0.30"), must(ParseSingle, "1.0.0.31"), must(ParseSingle, "1.0.0.32")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 31}),
		dst:     []string{"1.0.0.30", "1.0.0.32"},
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseSingle, "1.0.0.30"), must(ParseSingle, "1.0.0.31"), must(ParseSingle, "1.0.0.32")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 30}),
		dst:     []string{"1.0.0.31", "1.0.0.32"},
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{must(ParseSingle, "1.0.0.30"), must(ParseSingle, "1.0.0.31"), must(ParseSingle, "1.0.0.32")}},
		ip:      netip.AddrFrom4([4]byte{1, 0, 0, 32}),
		dst:     []string{"1.0.0.30", "1.0.0.31"},
	},
}

func TestExclude_DelAddress(t *testing.T) {
	for _, val := range testExcludeDelAddress {
		e := &Exclude{e: val.exclude, i: val.include}
		err := e.DelAddress(val.ip)
		switch {
		case err == nil && val.error == nil:
			if dst := e.Strings(); strings.Join(dst, "|") != strings.Join(val.dst, "|") {
				t.Errorf("expected value is %s, but got %s", val.dst, dst)
			}

		case err != nil && val.error != nil:
			if !errors.Is(err, val.error) {
				t.Errorf("expected value is %s, but got %s", val.error, err)
			}
		default:
			t.Errorf("expected value is %s, but got %s", val.error, err)
		}
	}
}

var testExcludeAddCIDR = []struct {
	include CIDR
	exclude Group
	val     Consecutive
	error   error
	dst     []string
}{
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{}},
		val:     must(ParseRange, "1.0.0.222-1.0.1.10"),
		error:   ErrNotInAddressRange,
	},
	{
		include: must(ParsePrefix, "1.0.0.0/24"),
		exclude: Group{arr: []Consecutive{}},
		val:     must(ParseRange, "1.0.0.222-1.0.0.240"),
		dst:     []string{"1.0.0.222-1.0.0.240"},
	},
	{
		include: must(ParsePrefix, "1.0.0.0/16"),
		exclude: Group{arr: []Consecutive{must(ParsePrefix, "1.0.0.0/24")}},
		val:     must(ParseRange, "1.0.0.12-1.0.0.22"),
		error:   ErrHasBeenExcluded,
	},
	{
		include: must(ParsePrefix, "1.0.0.0/16"),
		exclude: Group{arr: []Consecutive{must(ParsePrefix, "1.0.0.0/24")}},
		val:     must(ParseRange, "1.0.0.12-1.0.1.22"),
		error:   ErrHasBeenPartiallyExcluded,
	},
	{
		include: must(ParsePrefix, "1.0.0.0/16"),
		exclude: Group{arr: []Consecutive{must(ParsePrefix, "1.0.0.0/24")}},
		val:     must(ParseRange, "1.0.1.12-1.0.1.22"),
		dst:     []string{"1.0.0.0/24", "1.0.1.12-1.0.1.22"},
	},
}

func TestExclude_AddCIDR(t *testing.T) {
	for _, val := range testExcludeAddCIDR {
		e := &Exclude{e: val.exclude, i: val.include}
		err := e.AddCIDR(val.val)
		switch {
		case err == nil && val.error == nil:
			if dst := e.Strings(); strings.Join(dst, "|") != strings.Join(val.dst, "|") {
				t.Errorf("expected value is %s, but got %s", val.dst, dst)
			}

		case err != nil && val.error != nil:
			if !errors.Is(err, val.error) {
				t.Errorf("expected value is %s, but got %s", val.error, err)
			}
		default:
			t.Errorf("expected value is %s, but got %s", val.error, err)
		}
	}
}
