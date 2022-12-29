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
)

// 	case RequestURITooLong:
//		return "Request URI Too Long"

var ErrHijacked = errors.New("connection has been hijacked")
var ErrBadTrailer = errors.New("contain forbidden trailer")
