package pyrokinesis

import (
	"testing"
)

func TestBytes_ToString(t *testing.T) {
	Bytes.ToString([]byte("hello"))
}
