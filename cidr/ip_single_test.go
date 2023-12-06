package cidr

import (
	"math/big"
	"net/netip"
	"testing"
)

var testSingleLength = []struct {
	val    *Single
	expect *big.Int
}{
	{val: &Single{p: netip.AddrFrom4([4]byte{1, 1, 1, 1})}, expect: big.NewInt(1)},
	{
		val:    &Single{p: netip.AddrFrom16([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})},
		expect: big.NewInt(1),
	},
}

func TestSingle_Length(t *testing.T) {
	for _, val := range testSingleLength {
		b := val.val.Length()
		if b.Cmp(val.expect) != 0 {
			t.Errorf("expected value is %s, but got %s for Prefix length", val.expect.String(), b.String())
		}
	}
}

var testSingleFirstLastIP = []struct {
	val   *Single
	first netip.Addr
	last  netip.Addr
}{
	{
		val:   &Single{p: netip.AddrFrom4([4]byte{10, 0, 0, 0})},
		first: netip.AddrFrom4([4]byte{10, 0, 0, 0}),
		last:  netip.AddrFrom4([4]byte{10, 0, 0, 0}),
	},
}

func TestSingle_IP(t *testing.T) {
	for _, val := range testSingleFirstLastIP {
		if b := val.val.FirstIP(); b.Compare(val.first) != 0 {
			t.Errorf("expected value is %s, but got %s", val.first.String(), b.String())
		}

		if b := val.val.LastIP(); b.Compare(val.last) != 0 {
			t.Errorf("expected value is %s, but got %s", val.last.String(), b.String())
		}

	}
}

var testSingleNextIP = []struct {
	val  *Single
	ip   netip.Addr
	next netip.Addr
}{
	{
		val:  &Single{p: netip.AddrFrom4([4]byte{10, 0, 0, 0})},
		ip:   invalidIP,
		next: netip.AddrFrom4([4]byte{10, 0, 0, 0}),
	},
	{
		val:  &Single{p: netip.AddrFrom4([4]byte{10, 0, 0, 0})},
		ip:   netip.AddrFrom4([4]byte{10, 0, 0, 0}),
		next: invalidIP,
	},

	{
		val:  &Single{p: netip.AddrFrom16([16]byte{0: 1, 15: 5})},
		ip:   netip.AddrFrom16([16]byte{0: 1, 15: 5}),
		next: invalidIP,
	},
}

func TestSingle_NextIP(t *testing.T) {
	for _, val := range testSingleNextIP {
		if b := val.val.NextIP(val.ip); b.Compare(val.next) != 0 {
			t.Errorf("expected value is %s, but got %s", val.next.String(), b.String())
		}
	}
}

var testSingleString = []struct {
	src string
	dst string
}{
	{src: "1.0.0.0", dst: "1.0.0.0"},
	{src: "1.0.0.255", dst: "1.0.0.255"},
}

func TestSingle_String(t *testing.T) {
	for _, val := range testSingleString {
		p, _ := ParseSingle(val.src)
		if dst := p.String(); dst != val.dst {
			t.Errorf("expected value is %s, but got %s", val.dst, dst)
		}
	}
}
