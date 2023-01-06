package httpbody

import (
	"io"

	"go.x2ox.com/sorbifolia/http/httperr"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
)

type Memory struct {
	bufpool.Buffer
	p    int
	mode rwcMode
}

func (m *Memory) Read(p []byte) (n int, err error) {
	switch m.mode {
	case ModeReadWrite:
		m.mode = ModeRead
	case ModeRead:
	case ModeWrite:
		return 0, httperr.ErrNotYetReady
	case ModeClose:
		return 0, io.EOF
	default:
		panic("BUG: Memory should not exist in this state")
	}

	if m.Len() == m.p {
		return 0, io.EOF
	}
	n = copy(p, m.Buffer.B[m.p:])
	m.p += n
	return
}

func (m *Memory) Write(p []byte) (n int, err error) {
	switch m.mode {
	case ModeReadWrite:
		m.mode = ModeWrite
	case ModeWrite:
	case ModeRead, ModeClose:
		return 0, io.EOF
	default:
		panic("BUG: Memory should not exist in this state")
	}

	return m.Buffer.Write(p)
}

func (m *Memory) Close() error {
	switch m.mode {
	case ModeReadWrite:
		m.mode = ModeClose
	case ModeWrite, ModeRead:
		m.mode = ModeReadWrite
	case ModeClose:
	default:
		panic("BUG: Memory should not exist in this state")
	}
	return nil
}

func (m *Memory) Reset() {
	m.Buffer.Reset()
	m.p = 0
	m.mode = ModeReadWrite
}

func (m *Memory) release() { m.Reset(); _MemoryPool.Put(m) }
