package httpbody

import (
	"bytes"
	"io"
	"strconv"
	"sync"

	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
)

var (
	_ io.ReadCloser  = (*Chunked)(nil)
	_ io.WriteCloser = (*Chunked)(nil)
	_ HTTPBody       = (*Chunked)(nil)
)

type Chunked struct {
	Data, Header chan []byte

	m      rwcMode
	finish bool
	state  chunkedState
	once   sync.Once
	buf    bufpool.Buffer
}

func (c *Chunked) BodyReader() io.ReadCloser  { return c }
func (c *Chunked) BodyWriter() io.WriteCloser { return c }

func (c *Chunked) Write(p []byte) (n int, err error) {
	if !c.m.IsWrite() {
		return 0, io.EOF
	}

	pLen := len(p)

	for len(p) > 0 {
		if c.state == chunkedEND {
			break
		}

		if n, err = c.write(p); err != nil {
			return 0, err
		}
		p = p[n:]
	}

	if c.state == chunkedEND {
		err = io.EOF
	}

	return pLen, err
}

func (c *Chunked) write(p []byte) (n int, err error) {
	var (
		i   = bytes.Index(p, char.CRLF) // Key: Value\r\n\r\nBody
		buf = &c.buf
	)

	if i == -1 || i > 0 { // buf[\r], p[\n\r\n] -> i == 1
		if length := buf.Len(); p[0] == '\n' && length > 0 && buf.B[length-1] == '\r' { // Check "\r\n" is straddles the buffer.
			i = 0  // The data in buf is enough, no need to read again
			n = -1 // Two bytes will be discarded later
			buf.B = buf.B[:length-1]
		}
	}

	switch i {
	case 0:
	case -1: // TODO add size limit
		return buf.Write(p)
	default:
		n, _ = buf.Write(p[:i])
	}

	switch c.state {
	case chunkedRead:
		c.state++ // chunkedReadData. Has length, read data
		if buf.Len() == 1 && buf.B[0] == '0' {
			c.state++ // chunkedReadEnd. No length, read TrailerHeader or end
			close(c.Data)
		}
	case chunkedReadData:
		c.Data <- buf.Bytes()
		c.state--
	case chunkedReadEnd:
		if buf.Len() == 0 { // end
			c.state = chunkedEND
			close(c.Header)
		} else {
			c.Header <- buf.Bytes()
		}
	}

	n += 2 // Discard four bytes
	buf.Reset()

	return
}

func (c *Chunked) Read(p []byte) (n int, err error) {
	if !c.m.IsWrite() || (c.finish && c.buf.Len() == 0) {
		return 0, io.EOF
	}

	for {
		if len(p) == n {
			return
		}

		if c.buf.Len() > 0 {
			wn := copy(p[n:], c.buf.B)
			n += wn
			c.buf.Discard(0, wn)

			continue
		}

		if c.finish {
			return
		}

		data, ok := <-c.Data
		if ok {
			c.buf.B = strconv.AppendUint(c.buf.B, uint64(len(data)), 16)
			_, _ = c.buf.Write(char.CRLF)
			_, _ = c.buf.Write(data)
			continue
		}
		c.writeChunkedEnd()

		if c.Header != nil {
			data, ok = <-c.Header
		}
		if ok {
			_, _ = c.buf.Write(data)
			_, _ = c.buf.Write(char.CRLF)
		} else {
			_, _ = c.buf.Write(char.CRLF)
			c.finish = true
		}
	}
}

func (c *Chunked) Close() error { return nil }

func (c *Chunked) writeChunkedEnd() {
	c.once.Do(func() { _, _ = c.buf.Write([]byte("0\r\n")) })
}

type chunkedState uint8

const (
	chunkedRead chunkedState = iota
	chunkedReadData
	chunkedReadEnd
	chunkedEND
)