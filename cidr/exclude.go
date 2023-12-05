package cidr

import (
	"errors"
	"math/big"
	"net/netip"
)

type Exclude struct {
	e []Consecutive
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
	for _, v := range e.e {
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
	e.e = append(e.e, cidr)
	return nil
}

func (e *Exclude) DelAddress(addr netip.Addr) error {
	if !e.ContainsIP(addr) {
		return ErrNotInAddressRange
	}
	for i, v := range e.e {
		if !v.ContainsIP(addr) {
			continue
		}

		switch v.(type) {
		case Single:
			e.e = append(e.e[:i], e.e[i+1:]...)
			return nil

		case Range, Prefix:
			e.e[i] = Range{s: v.FirstIP(), e: addr.Prev()}
			e.e = append(e.e, Range{s: addr.Next(), e: v.LastIP()})
			return nil
		}
	}

	return nil
}

func (e *Exclude) AddAddress(addr netip.Addr) error {
	if !e.ContainsIP(addr) {
		return ErrNotInAddressRange
	}
	for _, v := range e.e {
		if v.ContainsIP(addr) {
			return ErrHasBeenExcluded
		}
	}
	e.e = append(e.e, Single{p: addr})
	return nil
}

func (e *Exclude) ContainsIP(addr netip.Addr) bool {
	if !e.i.ContainsIP(addr) {
		return false
	}
	for _, v := range e.e {
		if v.ContainsIP(addr) {
			return false
		}
	}
	return true
}

func (e *Exclude) NextIP(addr netip.Addr) netip.Addr {
	for {
		if addr = e.i.NextIP(addr); !addr.IsValid() {
			return addr
		}

		for _, v := range e.e {
			if v.ContainsIP(addr) {
				continue
			}
		}
		return addr
	}
}

func (e *Exclude) Length() *big.Int {
	bi := e.i.Length()
	for i := range e.e {
		bi.Sub(bi, e.e[i].Length())
	}
	return bi
}

func (e *Exclude) Contains(cidr CIDR) bool {
	if !e.i.Contains(cidr) {
		return false
	}

	for _, v := range e.e {
		if cidr.Contains(v) {
			return false
		}
	}
	return true
}
