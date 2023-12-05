package cidr

import (
	"math/big"
	"net/netip"
	"strings"
)

type Group struct {
	arr []Consecutive
}

func ParseGroup(s []string) (Group, error) {
	arr := make([]Consecutive, 0, len(s))
	for _, v := range s {
		switch {
		case strings.ContainsRune(v, '-'):
			val, err := ParseRange(v)
			if err != nil {
				return Group{}, err
			}
			arr = append(arr, val)

		case strings.ContainsRune(v, '/'):
			val, err := ParsePrefix(v)
			if err != nil {
				return Group{}, err
			}
			arr = append(arr, val)

		default:
			val, err := ParseSingle(v)
			if err != nil {
				return Group{}, err
			}
			arr = append(arr, val)
		}
	}

	return Group{arr}, nil
}

func (x Group) ContainsIP(ip netip.Addr) bool {
	for i := range x.arr {
		if x.arr[i].ContainsIP(ip) {
			return true
		}
	}
	return false
}

func (x Group) Length() *big.Int {
	bi := big.NewInt(0)
	for i := range x.arr {
		bi.Add(bi, x.arr[i].Length())
	}
	return bi
}

func (x Group) NextIP(ip netip.Addr) netip.Addr {
	if len(x.arr) == 0 {
		return netip.Addr{}
	}

	if !ip.IsValid() {
		return x.arr[0].FirstIP()
	}

	for i := range x.arr {
		if !x.arr[i].ContainsIP(ip) {
			continue
		}
		if addr := x.arr[i].NextIP(ip); addr.IsValid() {
			return addr
		}
		if len(x.arr)-1 != i { // traversal not finished
			return x.arr[i+1].FirstIP()
		}
		break
	}

	if x.arr[0].FirstIP().Is4() {
		return netip.IPv4Unspecified()
	}
	return netip.IPv6Unspecified()
}

func (x Group) Strings() []string {
	arr := make([]string, 0, len(x.arr))
	for _, v := range x.arr {
		arr = append(arr, v.String())
	}
	return arr
}

func (x Group) Contains(cidr CIDR) bool {
	for _, v := range x.arr {
		if v.Contains(cidr) {
			return false
		}
	}
	return true
}
