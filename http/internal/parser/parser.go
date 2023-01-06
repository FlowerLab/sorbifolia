package parser

import (
	"sync"
)

func AcquireRequestWriter() *RequestWriter {
	if v := _RequestWriterPool.Get(); v != nil {
		return v.(*RequestWriter)
	}
	return &RequestWriter{}
}

func ReleaseRequestWriter(r *RequestWriter) { r.Reset(); _RequestWriterPool.Put(r) }

var (
	_RequestWriterPool = sync.Pool{}
)
