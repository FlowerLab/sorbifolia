package httputils

import (
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()

	h := Get()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodGet) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestHead(t *testing.T) {
	t.Parallel()

	h := Head()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodHead) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestPost(t *testing.T) {
	t.Parallel()

	h := Post()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodPost) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestPut(t *testing.T) {
	t.Parallel()

	h := Put()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodPut) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestPatch(t *testing.T) {
	t.Parallel()

	h := Patch()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodPatch) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	h := Delete()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodDelete) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	h := Options()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodOptions) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestConnect(t *testing.T) {
	t.Parallel()

	h := Connect()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodConnect) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestTrace(t *testing.T) {
	t.Parallel()

	h := Trace()
	if rep, _, _ := h.test(); string(rep.Header.Method()) != string(MethodTrace) {
		t.Errorf("unexpected method: %s", rep.Header.Method())
	}
}

func TestNewMethod(t *testing.T) {
	t.Parallel()

	if string(NewMethod("Method")) != "Method" {
		t.Errorf("unexpected method")
	}
}
