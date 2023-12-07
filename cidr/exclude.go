package cidr

import (
	"math/big"
	"net/netip"
)

type Exclude struct {
	e Group
	i CIDR
}

func NewExclude(c CIDR, exclude ...Consecutive) (*Exclude, error) {
	e := &Exclude{i: c}
	for i := range exclude {
		if err := e.ExcludeCIDR(exclude[i]); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *Exclude) ExcludeCIDR(c Consecutive) error {
	switch e.i.Contains(c) {
	case ContainsPartially, ContainsNot:
		return ErrNotInAddressRange
	}

	for _, v := range e.e.arr {
		switch v.Contains(c) {
		case ContainsPartially:
			return ErrHasBeenPartiallyExcluded
		case Contains:
			return ErrHasBeenExcluded
		}
	}
	e.e.arr = append(e.e.arr, c)
	return nil
}

func (e *Exclude) DelExclude(addr netip.Addr) error {
	if !e.i.ContainsIP(addr) {
		return ErrNotInAddressRange
	}
	for i, v := range e.e.arr {
		if !v.ContainsIP(addr) {
			continue
		}

		switch v.(type) {
		case *Single:
			e.e.arr = append(e.e.arr[:i], e.e.arr[i+1:]...)
			return nil

		case *Range, *Prefix:
			vfi, vli, vl := v.FirstIP(), v.LastIP(), v.Length()

			switch {
			case vl.Cmp(big.NewInt(1)) == 0: // only one ip
				e.e.arr = append(e.e.arr[:i], e.e.arr[i+1:]...)

			case vfi.Compare(addr) == 0: // delete first ip
				e.e.arr[i] = &Range{s: vfi.Next(), e: vli}

			case vli.Compare(addr) == 0: // delete last ip
				e.e.arr[i] = &Range{s: vfi, e: vli.Prev()}

			default:
				e.e.arr[i] = &Range{s: v.FirstIP(), e: addr.Prev()}
				e.e.arr = append(e.e.arr, &Range{s: addr.Next(), e: v.LastIP()})
			}

			return nil
		}
	}

	return nil
}

func (e *Exclude) ExcludeAddress(addr netip.Addr) error {
	if !e.i.ContainsIP(addr) {
		return ErrNotInAddressRange
	}
	if e.e.ContainsIP(addr) {
		return ErrHasBeenExcluded
	}
	e.e.arr = append(e.e.arr, &Single{p: addr})
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
		if !e.e.ContainsIP(addr) {
			return addr
		}
	}
}

func (e *Exclude) Length() *big.Int { return big.NewInt(0).Sub(e.i.Length(), e.e.Length()) }

func (e *Exclude) Contains(c CIDR) ContainsStatus {
	switch e.i.Contains(c) {
	case ContainsPartially:
		return ContainsPartially

	case ContainsNot:
		return ContainsNot

	case Contains:
		switch e.e.Contains(c) {
		case ContainsPartially:
			return ContainsPartially
		case ContainsNot:
			return Contains
		case Contains:
			return ContainsNot
		}
	}

	return ContainsNot
}

func (e *Exclude) Strings() []string { return e.e.Strings() }
