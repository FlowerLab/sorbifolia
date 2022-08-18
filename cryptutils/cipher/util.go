package cipher

import (
	_ "crypto/cipher"
	_ "unsafe"
)

func maxUvarint(n int) uint64 {
	return 1<<uint(n*8) - 1
}

// put uint64 as big endian.
func putUvarint(bs []byte, u uint64) {
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(u >> uint(8*(len(bs)-1-i)))
	}
}

func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

//go:linkname xorBytes crypto/cipher.xorBytes
func xorBytes(dst, a, b []byte) int

// gfnDouble computes 2 * input in the field of 2^n elements.
// The irreducible polynomial in the finite field for n=128 is
// x^128 + x^7 + x^2 + x + 1 (equals 0x87)
// Constant-time execution in order to avoid side-channel attacks
func gfnDouble(input []byte) []byte {
	if len(input) != 16 {
		panic("Doubling in GFn only implemented for n = 128")
	}
	// If the first bit is zero, return 2L = L << 1
	// Else return (L << 1) xor 0^120 10000111
	shifted := shiftBytesLeft(input)
	shifted[15] ^= (input[0] >> 7) * 0x87
	return shifted
}

// shiftBytesLeft outputs the byte array corresponding to x << 1 in binary.
func shiftBytesLeft(x []byte) []byte {
	l := len(x)
	dst := make([]byte, l)
	for i := 0; i < l-1; i++ {
		dst[i] = (x[i] << 1) | (x[i+1] >> 7)
	}
	dst[l-1] = x[l-1] << 1
	return dst
}

// rightXor XORs smaller input (assumed Y) at the right of the larger input (assumed X)
func rightXor(X, Y []byte) []byte {
	offset := len(X) - len(Y)
	xored := make([]byte, len(X))
	copy(xored, X)
	for i := 0; i < len(Y); i++ {
		xored[offset+i] ^= Y[i]
	}
	return xored
}

// xorBytesMut assumes equal input length, replaces X with X XOR Y
func xorBytesMut(X, Y []byte) {
	for i := 0; i < len(X); i++ {
		X[i] ^= Y[i]
	}
}

// shiftNBytesLeft puts in dst the byte array corresponding to x << n in binary.
func shiftNBytesLeft(dst, x []byte, n int) {
	// Erase first n / 8 bytes
	copy(dst, x[n/8:])

	// Shift the remaining n % 8 bits
	bits := uint(n % 8)
	l := len(dst)
	for i := 0; i < l-1; i++ {
		dst[i] = (dst[i] << bits) | (dst[i+1] >> uint(8-bits))
	}
	dst[l-1] = dst[l-1] << bits

	// Append trailing zeroes
	dst = append(dst, make([]byte, n/8)...) //nolint:all
}
