package version

import (
	"testing"
)

func TestParseHttpVersion(t *testing.T) {
	major, minor, ok := parseHTTPVersion([]byte("HTTP/2.0"))
	if !ok {
		t.Error("parse error")
	}
	t.Log(major)
	t.Log(minor)
}
