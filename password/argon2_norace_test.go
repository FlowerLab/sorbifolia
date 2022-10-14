//go:build !race

package password

import (
	"crypto/rand"
	"errors"
	"testing"
	"time"
)

type errReader struct{}

func (e errReader) Read([]byte) (n int, err error) {
	return 0, errors.New("OEF")
}

func TestArgon2_MustGenerate(t *testing.T) {
	t.Parallel()

	time.Sleep(time.Second) // It changes rand.Reader
	g := New()

	rand.Reader = errReader{}

	defer func() { _ = recover() }()

	g.MustGenerate("123456")

	t.Error("fail")
}
