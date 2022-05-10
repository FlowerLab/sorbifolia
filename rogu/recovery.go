package rogu

import (
	"runtime"

	"go.uber.org/zap"
	"go.x2ox.com/sorbifolia/coarsetime"
)

// Recovery is the method used to catch panics.
func Recovery() func() {
	return func() {
		if err := recover(); err != nil {
			zap.L().Error("Recovery",
				zap.Time("time", coarsetime.FloorTime()),
				zap.Any("error", err),
				zap.String("stack", string(stack())),
			)
		}
	}
}

// Recover is a middleware that recovers from panics.
func Recover(fn func()) func() {
	return func() {
		defer Recovery()()
		fn()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return buf
}
