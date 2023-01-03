package util

import (
	"bufio"
	"io"
	"sync"
)

// bufio.NewReader()

// func Acquire(size ...int) *Buffer { return global.Acquire(size...) }
// func Release(b *Buffer)           { global.Release(b) }

var (
	_BufIOReaderPool = sync.Pool{New: func() any { return nil }}
	_BufIOWriterPool = sync.Pool{New: func() any { return nil }}
)

func AcquireBufIOReader(r io.Reader) *bufio.Reader {
	if v := _BufIOReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

func ReleaseBufIOReader(br *bufio.Reader) { _BufIOReaderPool.Put(br) }

func AcquireBufIOWriter(r io.Writer) *bufio.Writer {
	if v := _BufIOWriterPool.Get(); v != nil {
		br := v.(*bufio.Writer)
		br.Reset(r)
		return br
	}
	return bufio.NewWriter(r)
}

func ReleaseBufIOWriter(br *bufio.Writer) { _BufIOWriterPool.Put(br) }
