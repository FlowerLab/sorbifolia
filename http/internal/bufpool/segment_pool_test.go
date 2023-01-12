package bufpool

import (
	"testing"
)

func TestSegmentPool(t *testing.T) {
	for i := -1; i < max+1; i++ {
		sp := AcquireSegment(i)
		if cap(sp.B) == 0 {
			t.Error("Segment pool get failed")
		}
		sp.Release()
	}
}
