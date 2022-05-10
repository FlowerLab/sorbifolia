//go:build lumberjack

package rogu

import (
	"net/url"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	if err := ZapRegisterSink(new(LumberjackSink)); err != nil {
		panic(err)
	}
}

type LumberjackSink struct {
	*lumberjack.Logger
}

func (LumberjackSink) Sync() error                       { return nil }
func (LumberjackSink) String() string                    { return "lumberjack" }
func (l LumberjackSink) Sink(*url.URL) (zap.Sink, error) { return l, nil }

// Lumberjack create lumberjack sink instance
func Lumberjack(logPath string, MaxSize, MaxBackups, MaxAge int, compress bool) LumberjackSink {
	return LumberjackSink{Logger: &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    MaxSize,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		Compress:   compress,
	}}
}

// DefaultLumberjack return default config
func DefaultLumberjack() LumberjackSink {
	return Lumberjack("/var/log/rogu/log", 1024, 30, 90, true)
}
