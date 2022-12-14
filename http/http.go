package http

import (
	"crypto/tls"
	"io"
	"time"

	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/status"
	"go.x2ox.com/sorbifolia/http/version"
)

type Request struct {
	Method   method.Method
	Path     []byte
	Header   RequestHeader
	Body     io.ReadCloser
	Response Response
}

type Header struct {
	raw []KV

	Version       version.Version
	ContentLength int64
	ContentType   string
}

type RequestHeader struct {
	ContentLength int64
	Close         bool

	Accept         string
	AcceptEncoding string
	AcceptLanguage string
	UserAgent      string

	Host       []byte
	RemoteAddr []byte
	RequestURI []byte
	TLS        *tls.ConnectionState
}

type Response struct {
	Header ResponseHeader
	Body   io.ReadCloser
}

type ResponseHeader struct {
	Header

	StatusCode status.Status // e.g. 200

	Server []byte

	Date          time.Time
	ContentLength int64
	Close         bool
	Host          []byte
	RemoteAddr    []byte
	RequestURI    []byte
	TLS           *tls.ConnectionState
}

type URL struct {
	Scheme      []byte
	Opaque      []byte // encoded opaque data
	Host        []byte // host or host:port
	Path        []byte // path (relative paths may omit leading slash)
	RawPath     []byte // encoded path hint (see EscapedPath method)
	OmitHost    bool   // do not emit empty host (authority)
	ForceQuery  bool   // append a query ('?') even if RawQuery is empty
	RawQuery    []byte // encoded query values, without '?'
	Fragment    []byte // fragment for references, without '#'
	RawFragment []byte // encoded fragment hint (see EscapedFragment method)

	Username []byte
	Password *[]byte
}

type KV struct {
	K []byte
	V *[]byte
}
