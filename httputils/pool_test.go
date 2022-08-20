package httputils

import (
	"testing"
)

func TestGetRequestBuffer(t *testing.T) {
	r := getRequestBuffer()
	r.Put()
}

func TestGetRequestBufferPut(t *testing.T) {
	r := getRequestBuffer()
	r.Put()
}

func TestGetHTTPBuffer(t *testing.T) {
	r := getHttpBuffer()
	r.Put()
}

func TestGetHTTPBufferPut(t *testing.T) {
	r := getHttpBuffer()
	r.Put()
}

func TestPoolNil(t *testing.T) {
	t.Run("", func(t *testing.T) {
		r := getHttpBuffer()
		r = nil
		r.Put()
		n := getHttpBuffer()
		if n == nil {
			t.Error("err")
		}
	})
	t.Run("", func(t *testing.T) {
		r := getHttpBuffer()
		r.buf = nil
		r.Put()
		n := getHttpBuffer()
		if n.buf == nil {
			t.Error("err")
		}
	})
	t.Run("", func(t *testing.T) {
		r := getRequestBuffer()
		r = nil
		r.Put()
		n := getRequestBuffer()
		if n == nil {
			t.Error("err")
		}
	})
	t.Run("", func(t *testing.T) {
		r := getRequestBuffer()
		r.req = nil
		r.Put()
		n := getRequestBuffer()
		if n.req == nil {
			t.Error("err")
		}
	})
	t.Run("", func(t *testing.T) {
		r := getRequestBuffer()
		r.resp = nil
		r.Put()
		n := getRequestBuffer()
		if n.resp == nil {
			t.Error("err")
		}
	})
}
