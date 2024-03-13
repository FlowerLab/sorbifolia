//go:build !go1.22

package random

import (
	_ "unsafe"
)

//go:linkname _fastRand runtime.fastrand
func _fastRand() uint32

//go:linkname _fastRand64 runtime.fastrand64
func _fastRand64() uint64
