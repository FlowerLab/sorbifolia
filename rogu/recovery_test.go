package rogu

import (
	"testing"

	"go.uber.org/zap"
)

func TestRecovery(t *testing.T) {
	MustReplaceGlobals(DefaultZapConfig(DefaultZapEncoderConfig(), []string{"stdout"}, []string{"stderr"}))

	defer Recovery()()
	zap.L().Panic("test")
}

func TestRecover(t *testing.T) {
	MustReplaceGlobals(DefaultZapConfig(DefaultZapEncoderConfig(), []string{"stdout"}, []string{"stderr"}))

	Recover(func() {
		zap.L().Panic("test")
	})()
}
