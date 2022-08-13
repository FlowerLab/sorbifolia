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
