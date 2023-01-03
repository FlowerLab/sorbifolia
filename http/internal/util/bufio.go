package util

import (
	"bufio"
	"io"
	"sync"
)

var (
	_BRPool = sync.Pool{New: func() any { return nil }}
	_BWPool = sync.Pool{New: func() any { return nil }}
)

func AcquireBR(r io.Reader) *bufio.Reader {
	if v := _BRPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

func ReleaseBR(br *bufio.Reader) { _BRPool.Put(br) }

func AcquireBW(r io.Writer) *bufio.Writer {
	if v := _BWPool.Get(); v != nil {
		br := v.(*bufio.Writer)
		br.Reset(r)
		return br
	}
	return bufio.NewWriter(r)
}

func ReleaseBW(br *bufio.Writer) { _BWPool.Put(br) }
