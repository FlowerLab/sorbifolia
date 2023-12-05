package cidr

import (
	"fmt"
	"net/netip"
	"strings"
)

func ParseGroup(s []string) (*Group, error) {
	arr := make([]Consecutive, 0, len(s))
	for _, v := range s {
		val, err := ParseConsecutive(v)
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)
	}

	return &Group{arr}, nil
}

func ParseConsecutive(s string) (Consecutive, error) {
	s = strings.ReplaceAll(s, " ", "")
	switch {
	case strings.ContainsRune(s, '-'):
		return ParseRange(s)
	case strings.ContainsRune(s, '/'):
		return ParsePrefix(s)
	default:
		return ParseSingle(s)
	}
}

func ParseRange(s string) (*Range, error) {
	b, a, ok := strings.Cut(s, "-")
	if !ok {
		return nil, fmt.Errorf("cidr: parse range incorrect syntax, %s", s)
	}

	var (
		start, end netip.Addr
		err        error
	)
	if start, err = netip.ParseAddr(b); err != nil {
		return nil, err
	}
	if end, err = netip.ParseAddr(a); err != nil {
		return nil, err
	}

	return NewRange(start, end), nil
}

func ParsePrefix(s string) (*Prefix, error) {
	p, err := netip.ParsePrefix(s)
	if err != nil {
		return nil, err
	}
	return NewPrefix(p), nil
}

func ParseSingle(s string) (*Single, error) {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return nil, err
	}
	return NewSingle(addr), nil
}
