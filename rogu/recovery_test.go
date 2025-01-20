package rogu

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestRecovery(t *testing.T) {
	t.Parallel()

	MustReplaceGlobals(DefaultZapConfig(DefaultZapEncoderConfig(), []string{"stdout"}, []string{"stderr"}))

	defer Recovery()()
	zap.L().Panic("test")
}

func TestRecover(t *testing.T) {
	t.Parallel()

	MustReplaceGlobals(DefaultZapConfig(DefaultZapEncoderConfig(), []string{"stdout"}, []string{"stderr"}))

	Recover(func() {
		zap.L().Panic("test")
	})()

	fmt.Println("ASDasdasd")
	fmt.Println("\n 1 \n", string(stack()))
	fmt.Println("\n 2 \n", string(stack()))
	fmt.Println("\n 3 \n", string(stack()))
	fmt.Println("\n 4 \n", string(stack()))
	fmt.Println("\n 5 \n", string(stack()))
	fmt.Println("\n 6 \n", string(stack()))
	fmt.Println("\n 7 \n", string(stack()))
}
