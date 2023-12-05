package cidr

import (
	"math/big"
	"net/netip"
	"strings"
)

type Prefix struct {
	p netip.Prefix
}

func NewPrefix(p netip.Prefix) Prefix { return Prefix{p: p} }
func ParsePrefix(s string) (Prefix, error) {
	p, err := netip.ParsePrefix(strings.ReplaceAll(s, " ", ""))
	if err != nil {
		return Prefix{}, err
	}
	return NewPrefix(p), nil
}

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
	if x.p.Addr().Is4() {
		return netip.IPv4Unspecified()
	}
	return netip.IPv6Unspecified()
}

func (x Prefix) Contains(b CIDR) bool {
	if val, ok := b.(Consecutive); ok {
		return x.FirstIP().Compare(val.FirstIP()) <= 0 && x.LastIP().Compare(val.LastIP()) >= 0
	}

	return false
}
