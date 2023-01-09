package version

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	major, minor, ok := parseHTTPVersion([]byte("HTTP/2.0"))
	if !ok {
		t.Error("parse error")
	}
	t.Log(major)
	t.Log(minor)
}
