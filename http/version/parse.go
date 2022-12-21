package version

import (
	"bytes"
	"strconv"
)

var (
	httpVersion10     = []byte("HTTP/1.0")
	httpVersion11     = []byte("HTTP/1.1")
	httpVersionPrefix = []byte("HTTP/")
)

func Parse(ver []byte) (v Version, ok bool) {
	v.Major, v.Minor, ok = parseHTTPVersion(ver)
	return
}

func parseHTTPVersion(ver []byte) (major, minor int, ok bool) {
	switch {
	case bytes.Equal(ver, httpVersion10):
		return 1, 0, true
	case bytes.Equal(ver, httpVersion11):
		return 1, 1, true
	case !bytes.HasPrefix(ver, httpVersionPrefix):
		return 0, 0, false
	}

	length := len(ver)
	idx := bytes.IndexByte(ver[5:], '.')
	if idx == -1 {
		idx = length
	}

	maj, err := strconv.ParseUint(string(ver[5:idx]), 10, 0)
	if err != nil {
		return 0, 0, false
	}
	if length == idx {
		return int(maj), 0, true
	}

	var min uint64
	if min, err = strconv.ParseUint(string(ver[6+idx:]), 10, 0); err != nil {
		return 0, 0, false
	}
	return int(maj), int(min), true
}
