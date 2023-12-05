package cidr

import (
	"fmt"
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

func TestPrefix_FirstIP(t *testing.T) {
	for _, val := range testPrefixLength {
		fmt.Println(val.val.FirstIP(), val.val.LastIP())
	}
}
