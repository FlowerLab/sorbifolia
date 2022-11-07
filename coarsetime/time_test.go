package coarsetime

import (
	"testing"
)

func TestSince(t *testing.T) {
	t.Parallel()

	Since(Now())
}

func TestUntil(t *testing.T) {
	t.Parallel()

	Until(Now())
}

func TestNow(t *testing.T) {
	t.Parallel()

	Now()
}
