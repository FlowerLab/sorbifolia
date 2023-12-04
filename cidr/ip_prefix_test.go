package cidr

import (
	"fmt"
	"math/big"
	"net/netip"
	"testing"
)

var testPrefixLength = []struct {
	val    Prefix
	expect *big.Int
}{
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{0, 0, 0, 0}), 0)}, big.NewInt(4294967296)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 1)}, big.NewInt(2147483648)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 4)}, big.NewInt(268435456)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 8)}, big.NewInt(16777216)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 16)}, big.NewInt(65536)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 24)}, big.NewInt(256)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 0, 0, 0}), 32)}, big.NewInt(1)},

	{Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 0)}, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 80)}, big.NewInt(281474976710656)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 96)}, big.NewInt(4294967296)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 112)}, big.NewInt(65536)},
	{Prefix{p: netip.PrefixFrom(netip.AddrFrom16([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), 128)}, big.NewInt(1)},
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
