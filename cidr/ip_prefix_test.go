package cidr

import (
	"math/big"
	"net/netip"
	"testing"
)

var testPrefixLength = []struct {
	val    *Prefix
	expect *big.Int
}{
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{0, 0, 0, 0}), 0)}, expect: big.NewInt(4294967296)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 1)}, expect: big.NewInt(2147483648)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 4)}, expect: big.NewInt(268435456)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 8)}, expect: big.NewInt(16777216)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 16)}, expect: big.NewInt(65536)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 24)}, expect: big.NewInt(256)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 32)}, expect: big.NewInt(1)},

	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 0)}, expect: big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 80)}, expect: big.NewInt(281474976710656)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 96)}, expect: big.NewInt(4294967296)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 112)}, expect: big.NewInt(65536)},
	{val: &Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 128)}, expect: big.NewInt(1)},
}

func TestPrefix_Length(t *testing.T) {
	for _, val := range testPrefixLength {
		b := val.val.Length()
		if b.Cmp(val.expect) != 0 {
			t.Errorf("expected value is %s, but got %s for Prefix length", val.expect.String(), b.String())
		}
	}
}

var testPrefixFirstLastIP = []struct {
	val   *Prefix
	first netip.Addr
	last  netip.Addr
}{
	{
		val:   &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{10}), 24)},
		first: netip.AddrFrom4([4]byte{10, 0, 0, 0}),
		last:  netip.AddrFrom4([4]byte{10, 0, 0, 255}),
	},
}

func TestPrefix_IP(t *testing.T) {
	for _, val := range testPrefixFirstLastIP {
		if b := val.val.FirstIP(); b.Compare(val.first) != 0 {
			t.Errorf("expected value is %s, but got %s", val.first.String(), b.String())
		}

		if b := val.val.LastIP(); b.Compare(val.last) != 0 {
			t.Errorf("expected value is %s, but got %s", val.last.String(), b.String())
		}

	}
}

var testPrefixNextIP = []struct {
	val  *Prefix
	ip   netip.Addr
	next netip.Addr
}{
	{
		val:  &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{10}), 24)},
		ip:   invalidIP,
		next: netip.AddrFrom4([4]byte{10, 0, 0, 0}),
	},
	{
		val:  &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{10}), 24)},
		ip:   netip.AddrFrom4([4]byte{10, 0, 0, 0}),
		next: netip.AddrFrom4([4]byte{10, 0, 0, 1}),
	},
	{
		val:  &Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{10}), 24)},
		ip:   netip.AddrFrom4([4]byte{10, 0, 0, 255}),
		next: invalidIP,
	},

	{
		val:  &Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{0: 1, 15: 5}), 128)},
		ip:   netip.AddrFrom16([16]byte{0: 1, 15: 5}),
		next: invalidIP,
	},
}

func TestPrefix_NextIP(t *testing.T) {
	for _, val := range testPrefixNextIP {
		if b := val.val.NextIP(val.ip); b.Compare(val.next) != 0 {
			t.Errorf("expected value is %s, but got %s", val.next.String(), b.String())
		}
	}
}

var testPrefixString = []struct {
	src string
	dst string
}{
	{src: "1.0.0.0/32", dst: "1.0.0.0/32"},
	{src: "1.0.0.0/24", dst: "1.0.0.0/24"},
	{src: "0.0.0.0/0", dst: "0.0.0.0/0"},
}

func TestPrefix_String(t *testing.T) {
	for _, val := range testPrefixString {
		p, _ := ParsePrefix(val.src)
		if dst := p.String(); dst != val.dst {
			t.Errorf("expected value is %s, but got %s", val.dst, dst)
		}
	}
}
