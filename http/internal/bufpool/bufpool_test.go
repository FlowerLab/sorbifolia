package bufpool

import (
	"testing"
)

func TestAd(t *testing.T) {
	t.Log(-1 >> 1)
	t.Log(0 >> 1)
	t.Log(1230 >> 1)
	t.Log(64 >> 1)
}
