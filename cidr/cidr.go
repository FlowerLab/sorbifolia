package cidr

import (
	"math/big"
	"net/netip"
)

type CIDR interface {
	ContainsIP(netip.Addr) bool
	NextIP(netip.Addr) netip.Addr
	Length() *big.Int

	Contains(CIDR) ContainsStatus
}

type Consecutive interface {
	CIDR

	FirstIP() netip.Addr
	LastIP() netip.Addr
	String() string
}

type ContainsStatus uint8

const (
	Contains ContainsStatus = iota
	ContainsPartially
	ContainsNot
)

var (
	_ CIDR = (*Exclude)(nil)
	_ CIDR = (*Group)(nil)
	_ CIDR = (*Prefix)(nil)
	_ CIDR = (*Range)(nil)
	_ CIDR = (*Single)(nil)

	_ Consecutive = (*Prefix)(nil)
	_ Consecutive = (*Range)(nil)
	_ Consecutive = (*Single)(nil)
)
