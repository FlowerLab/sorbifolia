package httperr

import (
	"errors"
)

var (
	RequestURITooLong           = errors.New("request URI too long")
	RequestHeaderFieldsTooLarge = errors.New("request header fields too large")
	EntityTooLarge              = errors.New("entity too large")
	BodyTooLarge                = errors.New("body too large")
	BodyLengthMismatch          = errors.New("body length mismatch")
	ParseHTTPVersionErr         = errors.New("cannot find http version")
	ParseHTTPMethodErr          = errors.New("cannot find http method")
)

var (
	ErrHijacked   = errors.New("connection has been hijacked")
	ErrBadTrailer = errors.New("contain forbidden trailer")
)

var (
	ErrNotYetReady = errors.New("not yet ready")
)
