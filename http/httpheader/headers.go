package httpheader

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/netip"
	"strconv"
	"time"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type (
	Accept             []byte
	AcceptEncoding     []byte
	AcceptPatch        []byte
	AcceptPost         []byte
	AcceptRanges       []byte
	Allow              []byte
	Authorization      []byte
	ContentLanguage    []byte
	ContentEncoding    []byte
	ContentLength      []byte
	ContentLocation    []byte
	ContentRange       []byte
	AcceptLanguage     []byte
	ContentDisposition []byte
	ContentType        []byte
	Cookie             []byte
	Date               []byte
	Digest             []byte
	ETag               []byte
	Expires            []byte
	Host               []byte
	KeepAlive          []byte
	Referer            []byte
	Server             []byte
	TE                 []byte
	XForwardedHost     []byte
	XForwardedProto    []byte
	XFrameOptions      []byte
	Trailer            []byte
	TransferEncoding   []byte
	Upgrade            []byte
	UserAgent          []byte
	XForwardedFor      []byte
	SetCookie          []byte
	SetCookies         []SetCookie
	LastModified       []byte
	Location           []byte
	Origin             []byte
	Range              []byte
)

func (v Accept) Each(fn EachQualityValue)         { eachQualityValue(v, fn) }
func (v AcceptEncoding) Each(fn EachQualityValue) { eachQualityValue(v, fn) }
func (v AcceptRanges) Bytes() bool                { return bytes.Equal(v, char.Bytes) }
func (v AcceptRanges) None() bool                 { return len(v) == 0 || bytes.EqualFold(v, char.None) }
func (v AcceptLanguage) Each(fn EachQualityValue) { eachQualityValue(v, fn) }
func (v Allow) Each(fn EachValue)                 { eachValueWithComma(v, fn) }
func (v Authorization) Scheme() []byte            { return parseFirstValueOrNull(v, char.Space) }
func (v Authorization) Param() []byte             { return parseOtherValue(v, char.Space) }

func (v ContentLanguage) Each(fn EachValue) { eachValueWithComma(v, fn) }
func (v ContentLength) Length() (n int)     { return util.ToNonNegativeInt(v) }
func (v ContentRange) Unit() []byte         { return parseFirstValueOrNull(v, char.Space) }

func (v ContentRange) Start() int {
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

	return util.ToNonNegativeInt(b[:i])
}

func (v ContentRange) End() int {
	b := v
	i := bytes.IndexByte(b, char.Slash[0])
	if i < 0 {
		return -1
	}
	b = b[:i]

	if i = bytes.IndexByte(b, char.Hyphen); i < 0 {
		return -1
	}
	return util.ToNonNegativeInt(b[i+1:])
}

func (v ContentRange) Size() int {
	if i := bytes.IndexByte(v, char.Slash[0]); i >= 0 {
		return util.ToNonNegativeInt(v[i+1:])
	}
	return -1
}

func (v ContentDisposition) Type() []byte { return getHeaderValue(v) }

func (v ContentDisposition) Param(p []byte) (b []byte) {
	eachValueWithSemi(getHeaderParam(v), func(val []byte) bool {
		if k, value := parseKVWithEqual(val); bytes.EqualFold(k, p) {
			b = cleanQuotationMark(value)
			return false
		}
		return true
	})
	return
}

func (v ContentDisposition) Name() (b []byte)     { return v.Param([]byte("name")) }
func (v ContentDisposition) Filename() (b []byte) { return v.Param([]byte("filename")) }
func (v ContentType) MIME() []byte                { return parseFirstValue(v, char.Semi) }

func (v ContentType) Charset() []byte {
	var charset []byte = nil

	if i := bytes.Index(v, char.Charset); i >= 0 {
		charset = v[i+len(char.Charset):]
	}
	if len(charset) == 0 || charset[0] != char.Equal {
		return nil
	}
	charset = charset[1:]

	if i := bytes.IndexByte(charset, char.Semi); i >= 0 {
		charset = charset[:i]
	}

	return cleanQuotationMark(cleanSuffixSpace(charset))
}

func (v ContentType) Boundary() []byte {
	var boundary []byte = nil

	if i := bytes.Index(v, char.Boundary); i >= 0 {
		boundary = v[i+len(char.Boundary):]
	}
	if len(boundary) == 0 || boundary[0] != char.Equal {
		return nil
	}
	boundary = boundary[1:]

	if i := bytes.IndexByte(boundary, char.Semi); i >= 0 {
		boundary = boundary[:i]
	}

	return cleanSuffixSpace(boundary)
}

func (v Cookie) Each(fn EachKeyValue) {
	eachValueWithSemi(v, func(val []byte) bool { return fn(parseKVWithEqual(val)) })
}

func (v Date) Value() *time.Time { return toRFC1123(v) }

func (v Digest) Each(fn EachKeyValue) {
	eachValueWithComma(v, func(val []byte) bool {
		return fn(parseKVWithEqual(val))
	})
}

func (v ETag) Value() []byte {
	if bytes.HasPrefix(v, []byte("W/")) {
		return cleanQuotationMark(v[2:])
	}
	return cleanQuotationMark(v)
}

func (v Expires) Value() *time.Time { return toRFC1123(v) }

func (v Host) Host() []byte { return parseFirstValue(v, char.Colon) }

func (v Host) Port() uint16 {
	if i := bytes.IndexByte(v, char.Colon); i >= 0 {
		return uint16(util.ToNonNegativeInt(v[i+1:]))
	}
	return 0
}

func (v KeepAlive) Timeout() time.Duration {
	var d time.Duration = -1
	eachValueWithComma(v, func(value []byte) bool {
		if k, val := parseKVWithEqual(value); bytes.EqualFold(k, char.Timeout) {
			d = time.Duration(util.ToNonNegativeInt(val)) + time.Second
			return false
		}
		return true
	})
	return d
}
func (v KeepAlive) Max() int {
	i := -1
	eachValueWithComma(v, func(value []byte) bool {
		if k, val := parseKVWithEqual(value); bytes.EqualFold(k, char.Max) {
			i = util.ToNonNegativeInt(val)
			return false
		}
		return true
	})
	return i
}

func (v LastModified) Value() *time.Time { return toRFC1123(v) }

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
		return uint16(util.ToNonNegativeInt(b[i+1:]))
	}

	return 0
}

func (v Range) Unit() []byte {
	k, _ := parseKVWithEqual(v)
	return k
}

type Ranger struct {
	Start, End int // Valid range is 0 - 9223372036854775807, -1 means maximum or minimum
	Length     int // Range data Length
}

var ErrInvalidRangeSpecifier = errors.New("invalid ranges specifier")

func (r *Ranger) IsBad() bool {
	return (r.Start < 0 && r.End < 0 && r.Length < 0) || (r.Start > r.End)
}

func (r *Ranger) Reader(read io.ReadSeeker, length int) (io.Reader, error) {
	if r.IsBad() {
		return nil, ErrInvalidRangeSpecifier
	}

	if r.Length < 0 { // 100-
		if r.Start > length {
			return nil, ErrInvalidRangeSpecifier
		}
		r.Length = length - r.Start
		if _, err := read.Seek(int64(r.Start), io.SeekStart); err != nil {
			return nil, err
		}
		return read, nil
	}

	if r.Start < 0 && r.End < 0 { // -100
		if r.Length > length {
			return nil, ErrInvalidRangeSpecifier
		}
		if _, err := read.Seek(int64(length-r.Length), io.SeekStart); err != nil {
			return nil, err
		}
		return read, nil
	}

	if r.Start > length || r.Length > length { // 100-300
		return nil, ErrInvalidRangeSpecifier
	}
	if _, err := read.Seek(int64(r.Start), io.SeekStart); err != nil {
		return nil, err
	}
	return io.LimitReader(read, int64(r.Length)), nil
}

func (v Range) Each(fn EachRanger) {
	_, b := parseKVWithEqual(v)

	eachValueWithComma(b, func(value []byte) bool {
		i := bytes.IndexByte(value, char.Hyphen)

		switch {
		case i < 0:
			return fn(Ranger{Start: -1, End: -1, Length: -1})
		case i == 0: // -100
			return fn(Ranger{Start: -1, End: -1, Length: util.ToNonNegativeInt(value[1:])})
		case i == len(value)-1: // 100-
			return fn(Ranger{Start: util.ToNonNegativeInt(value[1:]), End: -1, Length: -1})
		}

		r := Ranger{Start: util.ToNonNegativeInt(value[:i]), End: util.ToNonNegativeInt(value[:i])}
		r.Length = r.End - r.Start
		return fn(r)
	})
}

func (v TE) Each(fn EachQualityValue)        { eachQualityValue(v, fn) }
func (v Trailer) Each(fn EachValue)          { eachValueWithComma(v, fn) }
func (v TransferEncoding) Each(fn EachValue) { eachValueWithComma(v, fn) }
func (v Upgrade) Each(fn EachValue)          { eachValueWithComma(v, fn) }
func (v XForwardedFor) Each(fn EachValue)    { eachValueWithComma(v, fn) }
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

func (v SetCookie) Cookie() (key, val []byte) {
	b := v
	if i := bytes.IndexByte(v, char.Semi); i >= 0 {
		b = v[:i]
	}
	i := bytes.IndexByte(v, char.Equal)

	return b[:i], cleanQuotationMark(b[i+1:])
}

func (v SetCookie) Expires() *time.Time {
	i := bytes.IndexByte(v, char.Semi)
	if i < 0 {
		return nil
	}
	b := v[i+1:]

	if i = bytes.Index(b, char.Expires); i < 0 || len(b) < i+2 || b[i+1] != char.Equal {
		return nil // Expires=
	}
	b = b[i+1:]

	if i = bytes.IndexByte(b, char.Semi); i >= 0 {
		b = b[:i]
	}

	b = cleanSuffixSpace(b)
	if len(b) == 0 {
		return nil
	}

	if t, err := time.Parse(time.RFC1123, pyrokinesis.Bytes.ToString(b)); err == nil {
		return &t
	}
	return nil
}

func (v SetCookie) MaxAge() int64 {
	i := bytes.IndexByte(v, char.Semi)
	if i < 0 {
		return -1
	}
	b := v[i+1:]

	if i = bytes.Index(b, char.MaxAge); i < 0 || len(b) < i+2 || b[i+1] != char.Equal {
		return -1 // Max-Age=
	}
	b = b[i+1:]

	if i = bytes.IndexByte(b, char.Semi); i >= 0 {
		b = b[:i]
	}

	b = cleanSuffixSpace(b)
	if len(b) == 0 {
		return -1
	}

	if n, err := strconv.ParseInt(pyrokinesis.Bytes.ToString(b), 10, 64); err == nil {
		return n
	}
	return -1
}

func (v SetCookie) Domain() []byte {
	i := bytes.IndexByte(v, char.Semi)
	if i < 0 {
		return nil
	}
	b := v[i+1:]

	if i = bytes.Index(b, char.Domain); i < 0 || len(b) < i+2 || b[i+1] != char.Equal {
		return nil // Domain=
	}
	b = b[i+1:]

	if i = bytes.IndexByte(b, char.Semi); i >= 0 {
		b = b[:i]
	}

	return cleanSuffixSpace(b)
}

func (v SetCookie) Path() []byte {
	i := bytes.IndexByte(v, char.Semi)
	if i < 0 {
		return nil
	}
	b := v[i+1:]

	if i = bytes.Index(b, char.Path); i < 0 || len(b) < i+2 || b[i+1] != char.Equal {
		return nil // Path=
	}
	b = b[i+1:]

	if i = bytes.IndexByte(b, char.Semi); i >= 0 {
		b = b[:i]
	}

	return cleanSuffixSpace(b)
}

func (v SetCookie) Secure() bool {
	if i := bytes.IndexByte(v, char.Semi); i >= 0 {
		return bytes.Contains(v[i+1:], char.Secure)
	}
	return false
}

func (v SetCookie) HttpOnly() bool {
	if i := bytes.IndexByte(v, char.Semi); i >= 0 {
		return bytes.Contains(v[i+1:], char.HttpOnly)
	}
	return false
}

// SameSite allows a server to define a cookie attribute making it impossible for
// the browser to send this cookie along with cross-site requests. The main
// goal is to mitigate the risk of cross-origin information leakage, and provide
// some protection against cross-site request forgery attacks.
//
// See https://tools.ietf.org/html/draft-ietf-httpbis-cookie-same-site-00 for details.
type SameSite uint8

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

func (v SetCookie) SameSite() SameSite {
	i := bytes.IndexByte(v, char.Semi)
	if i < 0 {
		return SameSiteDefaultMode
	}
	b := v[i+1:]

	if i = bytes.Index(b, char.SameSite); i < 0 || len(b) < i+2 || b[i+1] != char.Equal {
		return SameSiteDefaultMode // SameSite=
	}
	b = b[i+1:]

	if i = bytes.IndexByte(b, char.Semi); i >= 0 {
		b = b[:i]
	}
	b = cleanSuffixSpace(b)

	switch {
	case bytes.EqualFold(b, char.Lax):
		return SameSiteLaxMode
	case bytes.EqualFold(b, char.Strict):
		return SameSiteStrictMode
	case bytes.EqualFold(b, char.None):
		return SameSiteNoneMode
	}

	return SameSiteDefaultMode
}
