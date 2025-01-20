package bufferpool

import (
	"sync"

	"go.x2ox.com/sorbifolia/buffer"
)

var common sync.Pool

func Get() buffer.Buffer {
	b := common.Get()
	if b == nil {
		return &buffer.Byte{}
	}
	return b.(buffer.Buffer)
}

func Put(b buffer.Buffer) {
	if b != nil {
		b.Reset()
		common.Put(b)
	}
}
