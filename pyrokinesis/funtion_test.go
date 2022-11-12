package pyrokinesis

import (
	"testing"
)

type testCall struct{}

func (a testCall) A()  {}
func (a *testCall) B() {}

func TestCall(t *testing.T) {
	t.Parallel()

	s := &testCall{}
	Call(*s, "A", nil)
	Call(s, "B", nil)
}

func TestCallFail(t *testing.T) {
	t.Run("Test Call Fail cases", func(t *testing.T) {
		t.Parallel()

		defer func() { recover() }()

		Call(map[string]int{}, "A", nil)
		t.Error("Err")
	})

	t.Run("Test Call Fail cases", func(t *testing.T) {
		t.Parallel()

		defer func() { recover() }()

		Call(new(int), "A", nil)
		t.Error("Err")
	})

	t.Run("Test Call Fail cases", func(t *testing.T) {
		t.Parallel()

		defer func() { recover() }()

		Call(&testCall{}, "C", nil)
		t.Error("Err")
	})
}
