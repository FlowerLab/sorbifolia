package sa

import (
	"log/slog"
	"time"
)

// String returns an Attr for a string value.
func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

// Int64 returns an Attr for an int64.
func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

// Int converts an int to an int64 and returns
// an Attr with that value.
func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

// Uint64 returns an Attr for a uint64.
func Uint64(key string, v uint64) slog.Attr {
	return slog.Uint64(key, v)
}

// Float64 returns an Attr for a floating-point number.
func Float64(key string, v float64) slog.Attr {
	return slog.Float64(key, v)
}

// Bool returns an Attr for a bool.
func Bool(key string, v bool) slog.Attr {
	return slog.Bool(key, v)
}

// Time returns an Attr for a time.Time.
// It discards the monotonic portion.
func Time(key string, v time.Time) slog.Attr {
	return slog.Time(key, v)
}

// Duration returns an Attr for a time.Duration.
func Duration(key string, v time.Duration) slog.Attr {
	return slog.Duration(key, v)
}

// Group returns an Attr for a Group Value.
// The first argument is the key; the remaining arguments
// are converted to Attrs as in [Logger.Log].
//
// Use Group to collect several key-value pairs under a single
// key on a log line, or as the result of LogValue
// in order to log a single value as multiple Attrs.
func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}

// Any returns an Attr for the supplied value.
// See [AnyValue] for how values are treated.
func Any(key string, value any) slog.Attr {
	return slog.Attr{Key: key, Value: slog.AnyValue(value)}
}
