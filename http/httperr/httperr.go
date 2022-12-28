package httperr

import (
	"errors"
)

var (
	RequestHeaderFieldsTooLarge = errors.New("request header fields too large")
	EntityTooLarge              = errors.New("entity too large")
	BodyTooLarge                = errors.New("body too large")
	BodyLengthMismatch          = errors.New("body length mismatch")
)
