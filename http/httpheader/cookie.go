package httpheader

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/netip"
	"time"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type Cookie []byte

func (v Cookie) Each(fn EachKeyValue) {
	eachValueWithSemi(v, func(val []byte) bool {
		return fn(parseKVWithEqual(val))
	})
}

type Date []byte

func (v Date) Value() *time.Time { return toRFC1123(v) }

type Digest []byte

func (v Digest) Each(fn EachKeyValue) {
	eachValueWithComma(v, func(val []byte) bool {
		return fn(parseKVWithEqual(val))
	})
}

type ETag []byte

func (v ETag) Value() []byte {
	if bytes.HasPrefix(v, []byte("W/")) {
		return cleanQuotationMark(v[2:])
	}
	return cleanQuotationMark(v)
}

type Expires []byte

func (v Expires) Value() *time.Time { return toRFC1123(v) }

type Host []byte

func (v Host) Host() []byte {
	if i := bytes.IndexByte(v, char.Colon); i >= 0 {
		return v[:i]
	}
	return v
}

func (v Host) Port() uint16 {
	if i := bytes.IndexByte(v, char.Colon); i >= 0 {
		return uint16(toNonNegativeInt64(v[i+1:]))
	}
	return 0
}

type KeepAlive []byte

func (v KeepAlive) Timeout() time.Duration {
	var d time.Duration = -1
	eachValueWithComma(v, func(value []byte) bool {
		if k, val := parseKVWithEqual(value); bytes.EqualFold(k, char.Timeout) {
			d = time.Duration(toNonNegativeInt64(val)) + time.Second
			return false
		}
		return true
	})
	return d
}
func (v KeepAlive) Max() int64 {
	i := int64(-1)
	eachValueWithComma(v, func(value []byte) bool {
		if k, val := parseKVWithEqual(value); bytes.EqualFold(k, char.Max) {
			i = toNonNegativeInt64(val)
			return false
		}
		return true
	})
	return i
}

type LastModified []byte

func (v LastModified) Value() *time.Time { return toRFC1123(v) }

type Location []byte
type Origin []byte

func (v Origin) Null() bool {
	return len(v) == 0 || bytes.EqualFold(v, []byte("null"))
}

func (v Origin) Scheme() []byte {
	if i := bytes.Index(v, []byte("://")); i >= 0 {
		return v[:i]
	}
	return nil
}

func (v Origin) Hostname() []byte {
	i := bytes.Index(v, []byte("://"))
	if i < 0 {
		return nil
	}

	b := v[i+3:]
	if i = bytes.IndexByte(b, ':'); i >= 0 {
		return b[:i]
	}
	return b
}

func (v Origin) Port() uint16 {
	i := bytes.Index(v, []byte("://"))
	if i < 0 {
		return 0
	}

	b := v[i+3:]
	if i = bytes.IndexByte(b, ':'); i >= 0 {
		return uint16(toNonNegativeInt64(b[i+1:]))
	}

	return 0
}

type Range []byte

func (v Range) Unit() []byte {
	k, _ := parseKVWithEqual(v)
	return k
}

type Ranger struct {
	Start, End int64 // Valid range is 0 - 9223372036854775807, -1 means maximum or minimum
	Length     int64 // Range data Length
}

func (r *Ranger) IsBad() bool {
	return (r.Start < 0 && r.End < 0 && r.Length < 0) || (r.Start > r.End)
}

func (r *Ranger) Reader(read io.ReadSeeker, length int64) (io.Reader, error) {
	if r.IsBad() {
		return nil, errors.New("invalid ranges specifier")
	}

	if r.Length < 0 { // 100-
		if r.Start > length {
			return nil, errors.New("invalid ranges specifier")
		}
		r.Length = length - r.Start
		if _, err := read.Seek(r.Start, io.SeekStart); err != nil {
			return nil, err
		}
		return read, nil
	}

	if r.Start < 0 && r.End < 0 { // -100
		if r.Length > length {
			return nil, errors.New("invalid ranges specifier")
		}
		if _, err := read.Seek(length-r.Length, io.SeekStart); err != nil {
			return nil, err
		}
		return read, nil
	}

	if r.Start > length || r.Length > length { // 100-300
		return nil, errors.New("invalid ranges specifier")
	}
	if _, err := read.Seek(r.Start, io.SeekStart); err != nil {
		return nil, err
	}
	return io.LimitReader(read, r.Length), nil
}

func (v Range) Each(fn EachRanger) {
	_, b := parseKVWithEqual(v)

	eachValueWithComma(b, func(value []byte) bool {
		i := bytes.IndexByte(value, char.Hyphen)

		switch {
		case i < 0:
			return fn(Ranger{Start: -1, End: -1, Length: -1})
		case i == 0: // -100
			return fn(Ranger{Start: -1, End: -1, Length: toNonNegativeInt64(value[1:])})
		case i == len(value)-1: // 100-
			return fn(Ranger{Start: toNonNegativeInt64(value[1:]), End: -1, Length: -1})
		}

		r := Ranger{Start: toNonNegativeInt64(value[:i]), End: toNonNegativeInt64(value[:i])}
		r.Length = r.End - r.Start
		return fn(r)
	})
}

type Referer []byte
type Server []byte
type TE []byte

func (v TE) Each(fn EachQualityValue) { eachQualityValue(v, fn) }

type Trailer []byte

func (v Trailer) Each(fn EachValue) { eachValueWithComma(v, fn) }

type TransferEncoding []byte

func (v TransferEncoding) Each(fn EachValue) { eachValueWithComma(v, fn) }

type Upgrade []byte

func (v Upgrade) Each(fn EachValue) { eachValueWithComma(v, fn) }

type UserAgent []byte

type XForwardedFor []byte

func (v XForwardedFor) Each(fn EachValue) { eachValueWithComma(v, fn) }
func (v XForwardedFor) EachIP(fn func(addr net.IP) bool) {
	eachValueWithComma(v, func(val []byte) bool {
		return fn(net.ParseIP(pyrokinesis.Bytes.ToString(val)))
	})
}
func (v XForwardedFor) EachAddr(fn func(addr netip.Addr, err error) bool) {
	eachValueWithComma(v, func(val []byte) bool {
		return fn(netip.ParseAddr(pyrokinesis.Bytes.ToString(val)))
	})
}

type XForwardedHost []byte
type XForwardedProto []byte
type XFrameOptions []byte
