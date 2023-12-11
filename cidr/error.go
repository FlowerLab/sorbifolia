package cidr

import (
	"errors"
	"net/netip"
)

var (
	ErrNotInAddressRange        = errors.New("cidr: not in address range")
	ErrHasBeenExcluded          = errors.New("cidr: has been excluded")
	ErrHasBeenPartiallyExcluded = errors.New("cidr: has been partially excluded")
	ErrAddressRangeConflict     = errors.New("cidr: address range conflict")
)

var (
	invalidIP netip.Addr
)
