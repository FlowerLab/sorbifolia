package rogu

import (
	"log"
	"net/url"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestDefaultZapConfig(t *testing.T) {
	c := DefaultZapConfig(DefaultZapEncoderConfig(), []string{"stdout"}, []string{"stderr"})
	logger, err := c.Build()
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("test")
}

type TestSink struct {
}

func (ts *TestSink) String() string                  { return "TestSink" }
func (ts *TestSink) Sink(*url.URL) (zap.Sink, error) { return ts, nil }
func (ts *TestSink) Sync() error                     { return nil }
func (ts *TestSink) Close() error                    { return nil }
func (ts *TestSink) Write(p []byte) (n int, err error) {
	log.Print(string(p))
	return len(p), nil
}

func TestZapRegisterSink(t *testing.T) {
	if err := ZapRegisterSink(new(TestSink)); err != nil {
		t.Fatal(err)
	}
	if err := ZapRegisterSink(new(TestSink)); err == nil {
		t.Fatal("err")
	}
	c := DefaultZapConfig(DefaultZapEncoderConfig(),
		[]string{"TestSink:/log"},
		[]string{"stderr"})
	logger, err := c.Build()
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("test")
}

func TestMustReplaceGlobals(t *testing.T) {
	MustReplaceGlobals(DefaultZapConfig(DefaultZapEncoderConfig(),
		[]string{"stdout"},
		[]string{"stderr"}))
	zap.L().Info("test")

	Recover(func() {
		MustReplaceGlobals(zap.Config{
			Level: zap.AtomicLevel{},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:    "1",
				EncodeTime: nil,
			},
		})
	})
}
