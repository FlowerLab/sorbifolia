package cidr

import (
	"fmt"
	"math/big"
	"net/netip"
	"strings"
)

type Range struct {
	s, e netip.Addr
}

func NewRange(start, end netip.Addr) Range { return Range{s: start, e: end} }
func ParseRange(s string) (Range, error) {
	b, a, ok := strings.Cut(strings.ReplaceAll(s, " ", ""), "-")
	if !ok {
		return Range{}, fmt.Errorf("cidr: parse range incorrect syntax, %s", s)
	}

	var (
		start, end netip.Addr
		err        error
	)
	if start, err = netip.ParseAddr(b); err != nil {
		return Range{}, err
	}
	if end, err = netip.ParseAddr(a); err != nil {
		return Range{}, err
	}

	return NewRange(start, end), nil
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
func (x Range) Contains(b CIDR) bool {
	if val, ok := b.(Consecutive); ok {
		return x.FirstIP().Compare(val.FirstIP()) <= 0 && x.LastIP().Compare(val.LastIP()) >= 0
	}
	return false
}
