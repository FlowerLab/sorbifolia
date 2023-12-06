package cidr

import (
	"errors"
	"net/netip"
)

var (
	ErrNotInAddressRange        = errors.New("cidr: not in address range")
	ErrHasBeenExcluded          = errors.New("cidr: has been excluded")
	ErrHasBeenPartiallyExcluded = errors.New("cidr: has been partially excluded")
)

var (
	invalidIP netip.Addr
)
