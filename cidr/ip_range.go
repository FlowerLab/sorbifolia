package cidr

import (
	"fmt"
	"math/big"
	"net/netip"
)

type Range struct {
	s, e netip.Addr
}

func NewRange(start, end netip.Addr) Range { return Range{s: start, e: end} }

func (x Range) ContainsIP(ip netip.Addr) bool {
	return ip.Compare(x.s) >= 0 && ip.Compare(x.e) <= 0
}

func (x Range) Length() *big.Int {
	return big.NewInt(0).Sub(
		big.NewInt(0).SetBytes(x.s.AsSlice()),
		big.NewInt(0).SetBytes(x.e.AsSlice()),
	)
}

func (x Range) NextIP(ip netip.Addr) netip.Addr {
	if !ip.IsValid() {
		return x.s
	}
	if ip = ip.Next(); x.ContainsIP(ip) {
		return ip
	}

	if x.s.Is4() {
		return netip.IPv4Unspecified()
	}
	return netip.IPv6Unspecified()
}

func (x Range) FirstIP() netip.Addr { return x.s }
func (x Range) LastIP() netip.Addr  { return x.e }
func (x Range) String() string      { return fmt.Sprintf("%s-%s", x.s.String(), x.e.String()) }
func (x Range) Contains(c CIDR) ContainsStatus {
	if val, ok := c.(Consecutive); ok {
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
