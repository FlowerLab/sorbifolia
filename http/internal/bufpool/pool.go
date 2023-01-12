package bufpool

import (
	"sync"
)

var (
	readBufPool = sync.Pool{New: func() any { return &ReadBuffer{} }}
	bufPool     = sync.Pool{New: func() any { return &Buffer{} }}
)

func AcquireRBuf() *ReadBuffer { return readBufPool.Get().(*ReadBuffer) }
func AcquireBuf() *Buffer      { return bufPool.Get().(*Buffer) }
