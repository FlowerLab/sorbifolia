package bufpool

import (
	"sync"
)

var global = New()

func Acquire(size ...int) *Buffer { return global.Acquire(size...) }
func Release(b *Buffer)           { global.Release(b) }

type BufPool struct {
	p [steps]sync.Pool
}

func New() *BufPool {
	bp := &BufPool{}
	for i := range bp.p {
		bp.p[i].New = func() any { return &Buffer{B: make([]byte, 0, stepSize[i])} }
	}
	return bp
}

func (bp *BufPool) Release(b *Buffer) {
	if c := cap(b.B); c < min && c>>1 < max {
		b.B = b.B[:0]
		bp.p[indicator(cap(b.B))].Put(b)
	}
}

var stepSize = [steps]int{
	1 << (minStep + 0), 1 << (minStep + 1), 1 << (minStep + 2), 1 << (minStep + 3),
	1 << (minStep + 4), 1 << (minStep + 5), 1 << (minStep + 6), 1 << (minStep + 7),
	1 << (minStep + 8), 1 << (minStep + 9), 1 << (minStep + 10), 1 << (minStep + 11),
	1 << (minStep + 12), 1 << (minStep + 13), 1 << (minStep + 14), 1 << (minStep + 15),
	1 << (minStep + 16), 1 << (minStep + 17), 1 << (minStep + 18), 1 << (minStep + 19),
}

const (
	minStep = 6 // 2**6=64 is a CPU cache line size
	min     = 1 << (minStep + 0)
	max     = 1 << (minStep + 19)
	steps   = 20
)

func (bp *BufPool) Acquire(size ...int) *Buffer {
	var bs int
	if len(size) > 0 {
		bs = stepSize[indicator(size[0])]
	}
	return bp.p[bs].Get().(*Buffer)
}

func indicator(size int) int {
	switch {
	case size <= stepSize[0]:
		return 0
	case size&stepSize[18] >= stepSize[18]:
		return 19
	case size&stepSize[17] >= stepSize[17]:
		return 18
	case size&stepSize[16] >= stepSize[16]:
		return 17
	case size&stepSize[15] >= stepSize[15]:
		return 16
	case size&stepSize[14] >= stepSize[14]:
		return 15
	case size&stepSize[13] >= stepSize[13]:
		return 14
	case size&stepSize[12] >= stepSize[12]:
		return 13
	case size&stepSize[11] >= stepSize[11]:
		return 12
	case size&stepSize[10] >= stepSize[10]:
		return 11
	case size&stepSize[9] >= stepSize[9]:
		return 10
	case size&stepSize[8] >= stepSize[8]:
		return 9
	case size&stepSize[7] >= stepSize[7]:
		return 8
	case size&stepSize[6] >= stepSize[6]:
		return 7
	case size&stepSize[5] >= stepSize[5]:
		return 6
	case size&stepSize[4] >= stepSize[4]:
		return 5
	case size&stepSize[3] >= stepSize[3]:
		return 4
	case size&stepSize[2] >= stepSize[2]:
		return 3
	case size&stepSize[1] >= stepSize[1]:
		return 2
	case size&stepSize[0] > stepSize[0]:
		return 1
	default:
		panic("this is impossible")
	}
}
