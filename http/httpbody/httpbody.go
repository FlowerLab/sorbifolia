package httpbody

import (
	"io"
	"sync"
)

var (
	_ Pool = (*Chunked)(nil)
	_ Pool = (*Memory)(nil)
	_ Pool = (*TempFile)(nil)
	_ Pool = (*nobody)(nil)
)

type Pool interface {
	io.ReadWriteCloser
	Reset()
	release()
}

func Release(r Pool) { r.release() }

func AcquireMemory() *Memory {
	if v := _MemoryPool.Get(); v != nil {
		return v.(*Memory)
	}
	return &Memory{}
}

func AcquireChunked() *Chunked {
	if v := _ChunkedPool.Get(); v != nil {
		return v.(*Chunked)
	}
	return &Chunked{}
}

func AcquireTempFile() *TempFile {
	if v := _TempFilePool.Get(); v != nil {
		return v.(*TempFile)
	}
	return &TempFile{}
}

var (
	_MemoryPool   = sync.Pool{}
	_ChunkedPool  = sync.Pool{}
	_TempFilePool = sync.Pool{}
)
