package httpheader

import (
	"bytes"
	"strconv"
	"time"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type SetCookie []byte

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
		return bytes.Index(v[i+1:], char.Secure) >= 0
	}
	return false
}

func (v SetCookie) HttpOnly() bool {
	if i := bytes.IndexByte(v, char.Semi); i >= 0 {
		return bytes.Index(v[i+1:], char.HttpOnly) >= 0
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
