package rogu

import (
	"net/url"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RegisterSink interface {
	String() string
	Sink(*url.URL) (zap.Sink, error)
}

func ZapRegisterSink(rs ...RegisterSink) error {
	for _, v := range rs {
		if err := zap.RegisterSink(v.String(), v.Sink); err != nil {
			return err
		}
	}
	return nil
}

func DefaultZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func ZapConfig(level zapcore.Level, dev bool, encoding string, encoderConfig zapcore.EncoderConfig,
	outputPaths, errorOutputPaths []string,
) zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      dev,
		Encoding:         encoding, // json console
		EncoderConfig:    encoderConfig,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutputPaths,
	}
}

func DefaultZapConfig(ec zapcore.EncoderConfig, stdout, stderr []string) zap.Config {
	return ZapConfig(zapcore.DebugLevel, true, "console", ec, stdout, stderr)
}

func MustReplaceGlobals(config zap.Config) {
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
