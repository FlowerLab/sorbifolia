package cipher

import (
	"bytes"
	"crypto/cipher"
	"crypto/subtle"
	"errors"
	"math/bits"
)

type ocb struct {
	block     cipher.Block
	tagSize   int
	nonceSize int
	mask      mask
	// Optimized en/decrypt: For each nonce N used to en/decrypt, the 'Ktop'
	// internal variable can be reused for en/decrypting with nonces sharing
	// all but the last 6 bits with N. The prefix of the first nonce used to
	// compute the new Ktop, and the Ktop value itself, are stored in
	// reusableKtop. If using incremental nonces, this saves one block cipher
	// call every 63 out of 64 OCB encryptions, and stores one nonce and one
	// output of the block cipher in memory only.
	reusableKtop reusableKtop
}

type mask struct {
	// L_*, L_$, (L_i)_{i ∈ N}
	lAst []byte
	lDol []byte
	L    [][]byte
}

type reusableKtop struct {
	noncePrefix []byte
	Ktop        []byte
}

const (
	ocbDefaultTagSize   = 16
	ocbDefaultNonceSize = 15
)

const (
	enc = iota
	dec
)

func (o *ocb) NonceSize() int {
	return o.nonceSize
}

func (o *ocb) Overhead() int {
	return o.tagSize
}

// NewOCB returns an OCB instance with the given block cipher and default
// tag and nonce sizes.
func NewOCB(block cipher.Block) (cipher.AEAD, error) {
	return NewOCBWithNonceAndTagSize(block, ocbDefaultNonceSize, ocbDefaultTagSize)
}

// NewOCBWithNonceAndTagSize returns an OCB instance with the given block
// cipher, nonce length, and tag length. Panics on zero nonceSize and
// exceedingly long tag size.
//
// It is recommended to use at least 12 bytes as tag length.
func NewOCBWithNonceAndTagSize(
	block cipher.Block, nonceSize, tagSize int,
) (cipher.AEAD, error) {
	if block.BlockSize() != 16 {
		return nil, ocbError("Block cipher must have 128-bit blocks")
	}
	if nonceSize < 1 {
		return nil, ocbError("Incorrect nonce length")
	}
	if nonceSize >= block.BlockSize() {
		return nil, ocbError("Nonce length exceeds blocksize - 1")
	}
	if tagSize > block.BlockSize() {
		return nil, ocbError("Custom tag length exceeds blocksize")
	}
	return &ocb{
		block:     block,
		tagSize:   tagSize,
		nonceSize: nonceSize,
		mask:      initializeMaskTable(block),
		reusableKtop: reusableKtop{
			noncePrefix: nil,
			Ktop:        nil,
		},
	}, nil
}

func (o *ocb) Seal(dst, nonce, plaintext, adata []byte) []byte {
	if len(nonce) > o.nonceSize {
		panic("crypto/ocb: Incorrect nonce length given to OCB")
	}
	ret, out := sliceForAppend(dst, len(plaintext)+o.tagSize)
	o.crypt(enc, out, nonce, adata, plaintext)
	return ret
}

func (o *ocb) Open(dst, nonce, ciphertext, adata []byte) ([]byte, error) {
	if len(nonce) > o.nonceSize {
		panic("Nonce too long for this instance")
	}
	if len(ciphertext) < o.tagSize {
		return nil, ocbError("Ciphertext shorter than tag length")
	}
	sep := len(ciphertext) - o.tagSize
	ret, out := sliceForAppend(dst, len(ciphertext))
	ciphertextData := ciphertext[:sep]
	tag := ciphertext[sep:]
	o.crypt(dec, out, nonce, adata, ciphertextData)
	if subtle.ConstantTimeCompare(ret[sep:], tag) == 1 {
		ret = ret[:sep]
		return ret, nil
	}
	for i := range out {
		out[i] = 0
	}
	return nil, ocbError("Tag authentication failed")
}

// On instruction enc (resp. dec), crypt is the encrypt (resp. decrypt)
// function. It returns the resulting plain/ciphertext with the tag appended.
func (o *ocb) crypt(instruction int, Y, nonce, adata, X []byte) []byte {
	//
	// Consider X as a sequence of 128-bit blocks
	//
	// Note: For encryption (resp. decryption), X is the plaintext (resp., the
	// ciphertext without the tag).
	blockSize := o.block.BlockSize()

	//
	// Nonce-dependent and per-encryption variables
	//
	// Zero out the last 6 bits of the nonce into truncatedNonce to see if Ktop
	// is already computed.
	truncatedNonce := make([]byte, len(nonce))
	copy(truncatedNonce, nonce)
	truncatedNonce[len(truncatedNonce)-1] &= 192
	Ktop := make([]byte, blockSize) //nolint:staticcheck
	if bytes.Equal(truncatedNonce, o.reusableKtop.noncePrefix) {
		Ktop = o.reusableKtop.Ktop
	} else {
		// Nonce = num2str(TAGLEN mod 128, 7) || zeros(120 - bitlen(N)) || 1 || N
		paddedNonce := append(make([]byte, blockSize-1-len(nonce)), 1)
		paddedNonce = append(paddedNonce, truncatedNonce...)
		paddedNonce[0] |= byte(((8 * o.tagSize) % (8 * blockSize)) << 1)
		// Last 6 bits of paddedNonce are already zero. Encrypt into Ktop
		paddedNonce[blockSize-1] &= 192
		Ktop = paddedNonce
		o.block.Encrypt(Ktop, Ktop)
		o.reusableKtop.noncePrefix = truncatedNonce
		o.reusableKtop.Ktop = Ktop
	}

	// Stretch = Ktop || ((lower half of Ktop) XOR (lower half of Ktop << 8))
	xorHalves := make([]byte, blockSize/2)
	xorBytes(xorHalves, Ktop[:blockSize/2], Ktop[1:1+blockSize/2])
	stretch := append(Ktop, xorHalves...)
	bottom := int(nonce[len(nonce)-1] & 63)
	offset := make([]byte, len(stretch))
	shiftNBytesLeft(offset, stretch, bottom)
	offset = offset[:blockSize]

	//
	// Process any whole blocks
	//
	// Note: For encryption Y is ciphertext || tag, for decryption Y is
	// plaintext || tag.
	checksum := make([]byte, blockSize)
	m := len(X) / blockSize
	for i := 0; i < m; i++ {
		index := bits.TrailingZeros(uint(i + 1))
		if len(o.mask.L)-1 < index {
			o.mask.extendTable(index)
		}
		xorBytesMut(offset, o.mask.L[bits.TrailingZeros(uint(i+1))])
		blockX := X[i*blockSize : (i+1)*blockSize]
		blockY := Y[i*blockSize : (i+1)*blockSize]
		xorBytes(blockY, blockX, offset)
		switch instruction {
		case enc:
			o.block.Encrypt(blockY, blockY)
			xorBytesMut(blockY, offset)
			xorBytesMut(checksum, blockX)
		case dec:
			o.block.Decrypt(blockY, blockY)
			xorBytesMut(blockY, offset)
			xorBytesMut(checksum, blockY)
		}
	}
	//
	// Process any final partial block and compute raw tag
	//
	tag := make([]byte, blockSize)
	if len(X)%blockSize != 0 {
		xorBytesMut(offset, o.mask.lAst)
		pad := make([]byte, blockSize)
		o.block.Encrypt(pad, offset)
		chunkX := X[blockSize*m:]
		chunkY := Y[blockSize*m : len(X)]
		xorBytes(chunkY, chunkX, pad[:len(chunkX)])
		// P_* || bit(1) || zeroes(127) - len(P_*)
		switch instruction {
		case enc:
			paddedY := append(chunkX, byte(128))
			paddedY = append(paddedY, make([]byte, blockSize-len(chunkX)-1)...)
			xorBytesMut(checksum, paddedY)
		case dec:
			paddedX := append(chunkY, byte(128))
			paddedX = append(paddedX, make([]byte, blockSize-len(chunkY)-1)...)
			xorBytesMut(checksum, paddedX)
		}
		xorBytes(tag, checksum, offset)
		xorBytesMut(tag, o.mask.lDol)
		o.block.Encrypt(tag, tag)
		xorBytesMut(tag, o.hash(adata))
		copy(Y[blockSize*m+len(chunkY):], tag[:o.tagSize])
	} else {
		xorBytes(tag, checksum, offset)
		xorBytesMut(tag, o.mask.lDol)
		o.block.Encrypt(tag, tag)
		xorBytesMut(tag, o.hash(adata))
		copy(Y[blockSize*m:], tag[:o.tagSize])
	}
	return Y
}

// This hash function is used to compute the tag. Per design, on empty input it
// returns a slice of zeros, of the same length as the underlying block cipher
// block size.
func (o *ocb) hash(adata []byte) []byte {
	//
	// Consider A as a sequence of 128-bit blocks
	//
	A := make([]byte, len(adata))
	copy(A, adata)
	blockSize := o.block.BlockSize()

	//
	// Process any whole blocks
	//
	sum := make([]byte, blockSize)
	offset := make([]byte, blockSize)
	m := len(A) / blockSize
	for i := 0; i < m; i++ {
		chunk := A[blockSize*i : blockSize*(i+1)]
		index := bits.TrailingZeros(uint(i + 1))
		// If the mask table is too short
		if len(o.mask.L)-1 < index {
			o.mask.extendTable(index)
		}
		xorBytesMut(offset, o.mask.L[index])
		xorBytesMut(chunk, offset)
		o.block.Encrypt(chunk, chunk)
		xorBytesMut(sum, chunk)
	}

	//
	// Process any final partial block; compute final hash value
	//
	if len(A)%blockSize != 0 {
		xorBytesMut(offset, o.mask.lAst)
		// Pad block with 1 || 0 ^ 127 - bitlength(a)
		ending := make([]byte, blockSize-len(A)%blockSize)
		ending[0] = 0x80
		encrypted := append(A[blockSize*m:], ending...)
		xorBytesMut(encrypted, offset)
		o.block.Encrypt(encrypted, encrypted)
		xorBytesMut(sum, encrypted)
	}
	return sum
}

func initializeMaskTable(block cipher.Block) mask {
	//
	// Key-dependent variables
	//
	lAst := make([]byte, block.BlockSize())
	block.Encrypt(lAst, lAst)
	lDol := gfnDouble(lAst)
	L := make([][]byte, 1)
	L[0] = gfnDouble(lDol)

	return mask{
		lAst: lAst,
		lDol: lDol,
		L:    L,
	}
}

// Extends the L array of mask m up to L[limit], with L[i] = GfnDouble(L[i-1])
func (m *mask) extendTable(limit int) {
	for i := len(m.L); i <= limit; i++ {
		m.L = append(m.L, gfnDouble(m.L[i-1]))
	}
}

func ocbError(err string) error {
	return errors.New("crypto/ocb: " + err)
}
