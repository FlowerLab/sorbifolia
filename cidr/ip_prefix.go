package cidr

import (
	"math/big"
	"net/netip"
)

type Prefix struct {
	p netip.Prefix
}

func NewPrefix(p netip.Prefix) *Prefix { return &Prefix{p: p} }

func (x Prefix) ContainsIP(ip netip.Addr) bool { return x.p.Contains(ip) }

func (x Prefix) Length() *big.Int {
	return big.NewInt(0).Lsh(big.NewInt(1), uint(x.p.Addr().BitLen()-x.p.Bits()))
}

func (x Prefix) String() string      { return x.p.String() }
func (x Prefix) FirstIP() netip.Addr { return x.p.Addr() }
func (x Prefix) LastIP() netip.Addr {
	start := big.NewInt(0).SetBytes(x.p.Addr().AsSlice())
	lastIP := big.NewInt(0).Add(start, x.Length())
	lastIP = big.NewInt(0).Sub(lastIP, big.NewInt(1))

	addr, _ := netip.AddrFromSlice(lastIP.Bytes())
	return addr
}

func (x Prefix) NextIP(ip netip.Addr) netip.Addr {
	if !ip.IsValid() {
		return x.p.Addr()
	}
	if ip = ip.Next(); x.ContainsIP(ip) {
		return ip
	}

	return invalidIP
}

func (x Prefix) Contains(c CIDR) ContainsStatus {
	if val, ok := c.(Consecutive); ok {
		if x.LastIP().Compare(val.FirstIP()) < 0 || val.LastIP().Compare(x.FirstIP()) < 0 {
			return ContainsNot // x.start < x.end < c.start < c.end || c.start < c.end < x.start < x.end
		}

		xs := x.FirstIP().Compare(val.FirstIP()) <= 0
		xe := x.LastIP().Compare(val.LastIP()) >= 0
		if xs && xe { // x.start < c.start < c.end < x.end
			return Contains
		}
		if xs || xe { // x.start < c.start < x.end < c.end || c.start < x.start < c.end < x.end
			return ContainsPartially
		}
	}
	return ContainsNot
}
