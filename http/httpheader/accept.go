package httpheader

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
)

type Accept []byte

func (v Accept) Each(fn EachQualityValue) { eachQualityValue(v, fn) }

type AcceptEncoding []byte

func (v AcceptEncoding) Each(fn EachQualityValue) { eachQualityValue(v, fn) }

// AcceptPatch
//
// application/example, text/example
// text/example;charset=utf-8
// application/merge-patch+json
type AcceptPatch []byte
type AcceptPost []byte

// AcceptRanges: none

type AcceptRanges []byte

func (v AcceptRanges) Bytes() bool { return bytes.Equal(v, char.Bytes) }
func (v AcceptRanges) None() bool  { return len(v) == 0 || bytes.EqualFold(v, char.None) }

// Allow  GET, POST, HEAD
type Allow []byte

func (v Allow) Each(fn EachValue) { eachValueWithComma(v, fn) }

type Authorization []byte

func (v Authorization) Scheme() []byte {
	if i := bytes.IndexByte(v, char.Space); i >= 0 {
		return v[:i]
	}
	return nil
}

func (v Authorization) Param() []byte {
	if i := bytes.IndexByte(v, char.Space); i >= 0 {
		return v[i+1:]
	}
	return v
}

type ContentLanguage []byte

func (v ContentLanguage) Each(fn EachValue) { eachValueWithComma(v, fn) }

type ContentEncoding []byte
type ContentLength []byte

func (v ContentLength) Length() (n int64) {
	return util.ToNonNegativeInt64(v)
}

type ContentLocation []byte
type ContentRange []byte

func (v ContentRange) Unit() []byte {
	if i := bytes.IndexByte(v, char.Space); i >= 0 {
		return v[:i]
	}
	return nil
}

func (v ContentRange) Start() int64 {
	i := bytes.IndexByte(v, char.Space)
	if i < 0 {
		return -1
	}
	b := v[i+1:]

	if i = bytes.IndexByte(b, char.Slash[0]); i < 0 {
		return -1
	}
	b = b[:i]

	if i = bytes.IndexByte(b, char.Hyphen); i < 0 {
		return -1
	}

	return util.ToNonNegativeInt64(b[:i])
}

func (v ContentRange) End() int64 {
	i := bytes.IndexByte(v, char.Space)
	if i < 0 {
		return -1
	}
	b := v[i+1:]

	if i = bytes.IndexByte(b, char.Slash[0]); i < 0 {
		return -1
	}
	b = b[:i]

	if i = bytes.IndexByte(b, char.Hyphen); i < 0 {
		return -1
	}

	return util.ToNonNegativeInt64(b[i+1:])
}

func (v ContentRange) Size() int64 {
	i := bytes.IndexByte(v, char.Space)
	if i < 0 {
		return -1
	}
	b := v[i+1:]

	if i = bytes.IndexByte(b, char.Slash[0]); i < 0 {
		return -1
	}
	return util.ToNonNegativeInt64(b[i+1:])
}
