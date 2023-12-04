package cidr

import (
	"math/big"
	"net/netip"
)

type Single struct {
	p netip.Addr
}

func (x Single) ContainsIP(ip netip.Addr) bool { return x.p.Compare(ip) == 0 }

func (x Single) Length() *big.Int { return big.NewInt(1) }

func (x Single) NextIP(ip netip.Addr) netip.Addr {
	if !ip.IsValid() && x.p.Compare(ip) != 0 {
		return x.p
	}
	if x.p.Is4() {
		return netip.IPv4Unspecified()
	}
	return netip.IPv6Unspecified()
}

func (x Single) FirstIP() netip.Addr { return x.p }
func (x Single) LastIP() netip.Addr  { return x.p }

func (x Single) Contains(b CIDR) bool {
	s, ok := b.(Single)
	if !ok {
		return false
	}
	return s.p.Compare(x.p) == 0
}
