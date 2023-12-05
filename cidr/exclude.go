package cidr

import (
	"errors"
	"math/big"
	"net/netip"
)

type Exclude struct {
	e Group
	i CIDR
}

var (
	ErrNotInAddressRange        = errors.New("cidr: not in address range")
	ErrHasBeenExcluded          = errors.New("cidr: has been excluded")
	ErrHasBeenPartiallyExcluded = errors.New("cidr: has been partially excluded")
)

func (e *Exclude) AddCIDR(cidr Consecutive) error {
	if !e.Contains(cidr) {
		return ErrNotInAddressRange
	}
	for _, v := range e.e.arr {
		var (
			cs = v.ContainsIP(cidr.FirstIP()) // v.start < cidr.start < v.end
			ce = v.ContainsIP(cidr.LastIP())  // v.start < cidr.end < v.end
		)
		if cs && ce { // v.start < cidr.start < cidr.end < v.end
			return ErrHasBeenExcluded
		}
		if cs || ce { // v.start < cidr.start < v.end < cidr.end || cidr.start < v.start < cidr.end < v.end
			return ErrHasBeenPartiallyExcluded
		}
	}
	e.e.arr = append(e.e.arr, cidr)
	return nil
}

func (e *Exclude) DelAddress(addr netip.Addr) error {
	if !e.ContainsIP(addr) {
		return ErrNotInAddressRange
	}
	for i, v := range e.e.arr {
		if !v.ContainsIP(addr) {
			continue
		}

		switch v.(type) {
		case Single:
			e.e.arr = append(e.e.arr[:i], e.e.arr[i+1:]...)
			return nil

		case Range, Prefix:
			e.e.arr[i] = Range{s: v.FirstIP(), e: addr.Prev()}
			e.e.arr = append(e.e.arr, Range{s: addr.Next(), e: v.LastIP()})
			return nil
		}
	}

	return nil
}

func (e *Exclude) AddAddress(addr netip.Addr) error {
	if !e.ContainsIP(addr) {
		return ErrNotInAddressRange
	}
	for _, v := range e.e.arr {
		if v.ContainsIP(addr) {
			return ErrHasBeenExcluded
		}
	}
	e.e.arr = append(e.e.arr, Single{p: addr})
	return nil
}

func (e *Exclude) ContainsIP(addr netip.Addr) bool {
	return e.i.ContainsIP(addr) && !e.e.ContainsIP(addr)
}

func (e *Exclude) NextIP(addr netip.Addr) netip.Addr {
	for {
		if addr = e.i.NextIP(addr); !addr.IsValid() {
			return addr
		}

		var has bool
		for _, v := range e.e.arr {
			if has = v.ContainsIP(addr); has {
				break
			}
		}
		if !has {
			return addr
		}
	}
}

func (e *Exclude) Length() *big.Int { return big.NewInt(0).Sub(e.i.Length(), e.e.Length()) }

func (e *Exclude) Contains(cidr CIDR) bool {
	if !e.i.Contains(cidr) {
		return false
	}

	for _, v := range e.e.arr {
		if cidr.Contains(v) {
			return false
		}
	}
	return true
}
