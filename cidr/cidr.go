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
