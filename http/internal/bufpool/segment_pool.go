package bufpool

import (
	"sync"
)

var (
	global = NewSegmentPool()
)

func AcquireSegment(size int) *Buffer { return global.Acquire(size) }
func ReleaseSegment(b *Buffer)        { b.Reset(); global.Release(b) }

type SegmentPool struct {
	p [steps]sync.Pool
}

func NewSegmentPool() *SegmentPool {
	bp := &SegmentPool{}
	for i := range bp.p {
		bp.p[i].New = func() any { return &Buffer{B: make([]byte, 0, stepSize[i])} }
	}
	return bp
}

func (bp *SegmentPool) Release(b *Buffer) {
	if c := cap(b.B); c < min && c>>1 < max {
		b.B = b.B[:0]
		bp.p[indicator(cap(b.B))].Put(b)
	}
}

var stepSize = [steps]int{
	1 << (minStep + 0), 1 << (minStep + 1), 1 << (minStep + 2), 1 << (minStep + 3),
	1 << (minStep + 4), 1 << (minStep + 5), 1 << (minStep + 6), 1 << (minStep + 7),
	1 << (minStep + 8), 1 << (minStep + 9), 1 << (minStep + 10),
}

const (
	minStep = 6 // 2**6=64 is a CPU cache line size
	min     = 1 << (minStep + 0)
	max     = 1 << (minStep + 10)
	steps   = 11
)

func (bp *SegmentPool) Acquire(size int) *Buffer {
	return bp.p[indicator(size)].Get().(*Buffer)
}

func indicator(size int) int {
	switch {
	case size <= stepSize[0]:
		return 0
	case size >= stepSize[9]:
		return 10
	case size >= stepSize[8]:
		return 9
	case size >= stepSize[7]:
		return 8
	case size >= stepSize[6]:
		return 7
	case size >= stepSize[5]:
		return 6
	case size >= stepSize[4]:
		return 5
	case size >= stepSize[3]:
		return 4
	case size >= stepSize[2]:
		return 3
	case size >= stepSize[1]:
		return 2
	case size > stepSize[0]:
		return 1
	default:
		panic("BUG: this is impossible")
	}
}
