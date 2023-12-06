package cidr

import (
	"math/big"
	"net/netip"
	"testing"
)

var testRangeLength = []struct {
	val    *Range
	expect *big.Int
}{
	{val: &Range{s: netip.AddrFrom4([4]byte{1, 1, 1, 1}), e: netip.AddrFrom4([4]byte{1, 1, 1, 1})}, expect: big.NewInt(1)},
	{val: &Range{s: netip.AddrFrom4([4]byte{1, 1, 1, 1}), e: netip.AddrFrom4([4]byte{1, 1, 1, 2})}, expect: big.NewInt(2)},
	{val: &Range{s: netip.AddrFrom4([4]byte{}), e: netip.AddrFrom4([4]byte{255, 255, 255, 255})}, expect: big.NewInt(4294967296)},

	{
		val: &Range{
			s: netip.AddrFrom16([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			e: netip.AddrFrom16([16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}),
		},
		expect: big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil),
	},
	{
		val: &Range{
			s: netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
			e: netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
		},
		expect: big.NewInt(1),
	},
	{
		val: &Range{
			s: netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
			e: netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}),
		},
		expect: big.NewInt(2),
	},
}

func TestRange_Length(t *testing.T) {
	for _, val := range testRangeLength {
		b := val.val.Length()
		if b.Cmp(val.expect) != 0 {
			t.Errorf("expected value is %s, but got %s for Prefix length", val.expect.String(), b.String())
		}
	}
}

var testRangeFirstLastIP = []struct {
	val   *Range
	first netip.Addr
	last  netip.Addr
}{
	{
		val:   &Range{s: netip.AddrFrom4([4]byte{10, 0, 0, 0}), e: netip.AddrFrom4([4]byte{10, 0, 0, 255})},
		first: netip.AddrFrom4([4]byte{10, 0, 0, 0}),
		last:  netip.AddrFrom4([4]byte{10, 0, 0, 255}),
	},
}

func TestRange_IP(t *testing.T) {
	for _, val := range testRangeFirstLastIP {
		if b := val.val.FirstIP(); b.Compare(val.first) != 0 {
			t.Errorf("expected value is %s, but got %s", val.first.String(), b.String())
		}

		if b := val.val.LastIP(); b.Compare(val.last) != 0 {
			t.Errorf("expected value is %s, but got %s", val.last.String(), b.String())
		}

	}
}

var testRangeNextIP = []struct {
	val  *Range
	ip   netip.Addr
	next netip.Addr
}{
	{
		val:  &Range{s: netip.AddrFrom4([4]byte{10, 0, 0, 0}), e: netip.AddrFrom4([4]byte{10, 0, 0, 255})},
		ip:   netip.IPv4Unspecified(),
		next: netip.AddrFrom4([4]byte{10, 0, 0, 0}),
	},
	{
		val:  &Range{s: netip.AddrFrom4([4]byte{10, 0, 0, 0}), e: netip.AddrFrom4([4]byte{10, 0, 0, 255})},
		ip:   netip.AddrFrom4([4]byte{10, 0, 0, 0}),
		next: netip.AddrFrom4([4]byte{10, 0, 0, 1}),
	},
	{
		val:  &Range{s: netip.AddrFrom4([4]byte{10, 0, 0, 0}), e: netip.AddrFrom4([4]byte{10, 0, 0, 255})},
		ip:   netip.AddrFrom4([4]byte{10, 0, 0, 255}),
		next: netip.IPv4Unspecified(),
	},

	{
		val:  &Range{s: netip.AddrFrom16([16]byte{0: 1, 15: 5}), e: netip.AddrFrom16([16]byte{0: 1, 15: 5})},
		ip:   netip.AddrFrom16([16]byte{0: 1, 15: 5}),
		next: netip.IPv6Unspecified(),
	},
}

func TestRange_NextIP(t *testing.T) {
	for _, val := range testRangeNextIP {
		if b := val.val.NextIP(val.ip); b.Compare(val.next) != 0 {
			t.Errorf("expected value is %s, but got %s", val.next.String(), b.String())
		}
	}
}

var testRangeString = []struct {
	src string
	dst string
}{
	{src: "1.0.0.0-1.0.0.3", dst: "1.0.0.0-1.0.0.3"},
	{src: "1.0.0.0-1.0.0.255", dst: "1.0.0.0-1.0.0.255"},
	{src: "1.0.0.0-1.1.0.3", dst: "1.0.0.0-1.1.0.3"},
}

func TestRange_String(t *testing.T) {
	for _, val := range testRangeString {
		p, _ := ParseRange(val.src)
		if dst := p.String(); dst != val.dst {
			t.Errorf("expected value is %s, but got %s", val.dst, dst)
		}
	}
}
