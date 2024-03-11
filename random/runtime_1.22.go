//go:build go1.22

package random

import (
	_ "unsafe"
)

//go:linkname _fastRand runtime.cheaprand
func _fastRand() uint32

//go:linkname _fastRand64 runtime.cheaprand64
func _fastRand64() uint64
