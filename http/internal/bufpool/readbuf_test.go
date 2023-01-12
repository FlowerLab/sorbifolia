package bufpool

import (
	"fmt"
	"testing"
	"time"
)

func TestReadBufferAcquireReleaseSerial(t *testing.T) {
	testReadBufferAcquireRelease(t)
}

func TestReadBufferAcquireReleaseConcurrent(t *testing.T) {
	concurrency := 10
	ch := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			testReadBufferAcquireRelease(t)
			ch <- struct{}{}
		}()
	}

	for i := 0; i < concurrency; i++ {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatalf("timeout!")
		}
	}
}

func testReadBufferAcquireRelease(t *testing.T) {
	for i := 0; i < 10; i++ {
		expectedS := fmt.Sprintf("num %d", i)
		b := AcquireRBuf()
		b.B = append(b.B, "num "...)
		b.B = append(b.B, fmt.Sprintf("%d", i)...)
		if string(b.B) != expectedS {
			t.Fatalf("unexpected result: %q. Expecting %q", b.B, expectedS)
		}
		b.Release()
	}
}
