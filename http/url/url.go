package url

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
)

type URL struct {
	Scheme   []byte
	Host     []byte // host or host:port
	Path     []byte // path (relative paths may omit leading slash)
	Query    []byte // encoded query values, without '?'
	Fragment []byte // fragment for references, without '#'

	Username []byte
	Password *[]byte
}

// SetSchemeBytes sets URI scheme, i.e. http, https, ftp, etc.
func (u *URL) SetSchemeBytes(scheme []byte) {
	u.Scheme = util.ToLower(scheme)
}

func (u *URL) Parse(host, path []byte, isTLS bool) error {
	if len(host) == 0 || bytes.Contains(path, char.ColonSlashSlash) {
		_, host, path = splitHostURI(host, path)
	}
	if isTLS {
		u.SetSchemeBytes(char.HTTPS)
	} else {
		u.SetSchemeBytes(char.HTTP)
	}

	if n := bytes.IndexByte(host, char.At); n >= 0 {
		auth := host[:n]
		host = host[n+1:]

		if n = bytes.IndexByte(auth, char.Colon); n >= 0 {
			u.Username = auth[:n]
			pwd := auth[n+1:]
			u.Password = &pwd
		} else {
			u.Username = auth
		}
	}

	if parsedHost, err := parseHost(host); err != nil {
		return err
	} else {
		u.Host = parsedHost
	}

	queryIndex := bytes.IndexByte(path, char.QuestionMark)
	fragmentIndex := bytes.IndexByte(path, char.Hashtag)
	// Ignore query in fragment part
	if fragmentIndex >= 0 && queryIndex > fragmentIndex {
		queryIndex = -1
	}

	if queryIndex < 0 && fragmentIndex < 0 {
		u.Path = path
		return nil
	}

	if queryIndex >= 0 {
		// Path is everything up to the start of the query
		u.Path = path[:queryIndex]

		if fragmentIndex < 0 {
			u.Query = path[queryIndex+1:]
		} else {
			u.Query = path[queryIndex+1 : fragmentIndex]
			u.Fragment = path[fragmentIndex+1:]
		}
		return nil
	}

	u.Fragment = path[fragmentIndex+1:]
	return nil
}

func splitHostURI(host, uri []byte) ([]byte, []byte, []byte) {
	n := bytes.Index(uri, char.SlashSlash)
	if n < 0 {
		return char.HTTP, host, uri
	}
	scheme := uri[:n]
	if bytes.IndexByte(scheme, '/') >= 0 {
		return char.HTTP, host, uri
	}
	if len(scheme) > 0 && scheme[len(scheme)-1] == ':' {
		scheme = scheme[:len(scheme)-1]
	}
	n += len(char.SlashSlash)
	uri = uri[n:]
	n = bytes.IndexByte(uri, '/')
	nq := bytes.IndexByte(uri, '?')
	if nq >= 0 && nq < n {
		// A hack for urls like foobar.com?a=b/xyz
		n = nq
	} else if n < 0 {
		// A hack for bogus urls like foobar.com?a=b without
		// slash after host.
		if nq >= 0 {
			return scheme, uri[:nq], uri[nq:]
		}
		return scheme, uri, char.Slash
	}
	return scheme, uri[:n], uri[n:]
}

// parseHost parses host as an authority without user
// information. That is, as host[:port].
//
// Based on https://github.com/golang/go/blob/8ac5cbe05d61df0a7a7c9a38ff33305d4dcfea32/src/net/url/url.go#L619
//
// The host is parsed and unescaped in place overwriting the contents of the host parameter.
func parseHost(host []byte) ([]byte, error) {
	if len(host) > 0 && host[0] == '[' {
		// Parse an IP-Literal in RFC 3986 and RFC 6874.
		// E.g., "[fe80::1]", "[fe80::1%25en0]", "[fe80::1]:80".
		i := bytes.LastIndexByte(host, ']')
		if i < 0 {
			return nil, errors.New("missing ']' in host")
		}
		colonPort := host[i+1:]
		if !validOptionalPort(colonPort) {
			return nil, fmt.Errorf("invalid port %q after host", colonPort)
		}

		// RFC 6874 defines that %25 (%-encoded percent) introduces
		// the zone identifier, and the zone identifier can use basically
		// any %-encoding it likes. That's different from the host, which
		// can only %-encode non-ASCII bytes.
		// We do impose some restrictions on the zone, to avoid stupidity
		// like newlines.
		zone := bytes.Index(host[:i], []byte("%25"))
		if zone >= 0 {
			host1, err := unescape(host[:zone], encodeHost)
			if err != nil {
				return nil, err
			}
			host2, err := unescape(host[zone:i], encodeZone)
			if err != nil {
				return nil, err
			}
			host3, err := unescape(host[i:], encodeHost)
			if err != nil {
				return nil, err
			}
			return append(host1, append(host2, host3...)...), nil
		}
	} else if i := bytes.LastIndexByte(host, ':'); i != -1 {
		colonPort := host[i:]
		if !validOptionalPort(colonPort) {
			return nil, fmt.Errorf("invalid port %q after host", colonPort)
		}
	}

	var err error
	if host, err = unescape(host, encodeHost); err != nil {
		return nil, err
	}
	return host, nil
}

// validOptionalPort reports whether port is either an empty string
// or matches /^:\d*$/
func validOptionalPort(port []byte) bool {
	if len(port) == 0 {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

type EscapeError string

func (e EscapeError) Error() string {
	return "invalid URL escape " + strconv.Quote(string(e))
}

type InvalidHostError string

func (e InvalidHostError) Error() string {
	return "invalid character " + strconv.Quote(string(e)) + " in host name"
}

// unescape unescapes a string; the mode specifies
// which section of the URL string is being unescaped.
//
// Based on https://github.com/golang/go/blob/8ac5cbe05d61df0a7a7c9a38ff33305d4dcfea32/src/net/url/url.go#L199
//
// Unescapes in place overwriting the contents of s and returning it.
func unescape(s []byte, mode encoding) ([]byte, error) {
	// Count %, check that they're well-formed.
	n := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++
			if i+2 >= len(s) || !isHEX(s[i+1]) || !isHEX(s[i+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[:3]
				}
				return nil, EscapeError(s)
			}
			// Per https://tools.ietf.org/html/rfc3986#page-21
			// in the host component %-encoding can only be used
			// for non-ASCII bytes.
			// But https://tools.ietf.org/html/rfc6874#section-2
			// introduces %25 being allowed to escape a percent sign
			// in IPv6 scoped-address literals. Yay.
			if mode == encodeHost && unHEX(s[i+1]) < 8 && !bytes.Equal(s[i:i+3], []byte("%25")) {
				return nil, EscapeError(s[i : i+3])
			}
			if mode == encodeZone {
				// RFC 6874 says basically "anything goes" for zone identifiers
				// and that even non-ASCII can be redundantly escaped,
				// but it seems prudent to restrict %-escaped bytes here to those
				// that are valid host name bytes in their unescaped form.
				// That is, you can use escaping in the zone identifier but not
				// to introduce bytes you couldn't just write directly.
				// But Windows puts spaces here! Yay.
				v := unHEX(s[i+1])<<4 | unHEX(s[i+2])
				if !bytes.Equal(s[i:i+3], []byte("%25")) && v != ' ' && shouldEscape(v, encodeHost) {
					return nil, EscapeError(s[i : i+3])
				}
			}
			i += 3
		default:
			if (mode == encodeHost || mode == encodeZone) && s[i] < 0x80 && shouldEscape(s[i], mode) {
				return nil, InvalidHostError(s[i : i+1])
			}
			i++
		}
	}

	if n == 0 {
		return s, nil
	}

	t := s[:0]
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '%':
			t = append(t, unHEX(s[i+1])<<4|unHEX(s[i+2]))
			i += 2
		default:
			t = append(t, s[i])
		}
	}
	return t, nil
}

type encoding int

const (
	encodeHost encoding = 1 + iota
	encodeZone
)

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
//
// Please be informed that for now shouldEscape does not check all
// reserved characters correctly. See golang.org/issue/5684.
//
// Based on https://github.com/golang/go/blob/8ac5cbe05d61df0a7a7c9a38ff33305d4dcfea32/src/net/url/url.go#L100
func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == encodeHost || mode == encodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	if c == '-' || c == '_' || c == '.' || c == '~' { // §2.3 Unreserved characters (mark)
		return false
	}

	// Everything else must be escaped.
	return true
}

func isHEX(c byte) bool {
	return ('0' <= c && c <= '9') ||
		('a' <= c && c <= 'f') ||
		('A' <= c && c <= 'F')
}

func unHEX(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
