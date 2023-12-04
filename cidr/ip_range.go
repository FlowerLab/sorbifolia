package cidr

import (
	"math/big"
	"net/netip"
)

type Range struct {
	s, e netip.Addr
}

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
	if !ip.IsValid() { // 第一次调用
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
func (x Range) Contains(b CIDR) bool {
	if val, ok := b.(Consecutive); ok {
		return x.FirstIP().Compare(val.FirstIP()) <= 0 && x.LastIP().Compare(val.LastIP()) >= 0
	}
	return false
}
