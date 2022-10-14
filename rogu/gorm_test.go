//go:build gorm

package rogu

import (
	"testing"
)

func TestGorm(t *testing.T) {
	t.Parallel()

	MustReplaceGlobals(DefaultZapConfig(DefaultZapEncoderConfig(), []string{"stdout"}, []string{"stderr"}))

	g := Gorm(nil)
	g.Error(nil, "Asd")
}
