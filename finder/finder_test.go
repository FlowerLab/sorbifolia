package finder

import (
	"testing"
)

func TestAsync(t *testing.T) {
	t.Parallel()

	fn := func() string { return "async" }
	ch := Async(fn)
	s := <-ch
	t.Log(s)
}

func TestAsyncOR(t *testing.T) {
	t.Parallel()

	fn := func() string { return "async" }
	ch := AsyncOR(fn)
	s := <-ch
	t.Log(s)
}

func TestAsyncC(t *testing.T) {
	t.Parallel()

	fn := func() string { return "async" }
	ch := AsyncC(fn)
	s := <-ch
	t.Log(s)
}
