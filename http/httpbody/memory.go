package httpbody

import (
	"io"

	"go.x2ox.com/sorbifolia/http/internal/bufpool"
)

var (
	_ io.ReadCloser  = (*Memory)(nil)
	_ io.WriteCloser = (*Memory)(nil)
	_ HTTPBody       = (*Memory)(nil)
)

type Memory struct {
	bufpool.Buffer
	p    int
	mode rwcMode
}

func (m *Memory) Read(p []byte) (n int, err error) {
	if m.Len() == m.p {
		return 0, io.EOF
	}
	n = copy(p, m.Buffer.B[m.p:])
	m.p += n
	return
}

func (m *Memory) Close() error               { return nil }
func (m *Memory) BodyReader() io.ReadCloser  { return m }
func (m *Memory) BodyWriter() io.WriteCloser { return m }
