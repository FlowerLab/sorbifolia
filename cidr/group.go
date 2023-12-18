package cidr

import (
	"math/big"
	"net/netip"
)

type Group struct {
	Arr []Consecutive
}

func (x *Group) ContainsIP(ip netip.Addr) bool {
	for i := range x.Arr {
		if x.Arr[i].ContainsIP(ip) {
			return true
		}
	}
	return false
}

func (x *Group) Length() *big.Int {
	bi := big.NewInt(0)
	for i := range x.Arr {
		bi.Add(bi, x.Arr[i].Length())
	}
	return bi
}

func (x *Group) NextIP(ip netip.Addr) netip.Addr {
	if len(x.Arr) == 0 {
		return invalidIP
	}

	if !ip.IsValid() {
		return x.Arr[0].FirstIP()
	}

	for i := range x.Arr {
		if !x.Arr[i].ContainsIP(ip) {
			continue
		}
		if addr := x.Arr[i].NextIP(ip); addr.IsValid() {
			return addr
		}
		if len(x.Arr)-1 != i { // traversal not finished
			return x.Arr[i+1].FirstIP()
		}
		break
	}

	return invalidIP
}

func (x *Group) Strings() []string {
	arr := make([]string, 0, len(x.Arr))
	for _, v := range x.Arr {
		arr = append(arr, v.String())
	}
	return arr
}

func (x *Group) Contains(cidr CIDR) ContainsStatus {
	for _, v := range x.Arr {
		switch v.Contains(cidr) {
		case ContainsPartially: // TODO: dealing with splits
			return ContainsPartially
		case ContainsNot:
			continue
		case Contains:
			return Contains
		}
	}
	return ContainsNot
}

func (x *Group) AddCIDR(c Consecutive) error {
	switch x.Contains(c) {
	case ContainsPartially, Contains:
		return ErrAddressRangeConflict
	}
	x.Arr = append(x.Arr, c)
	return nil
}
