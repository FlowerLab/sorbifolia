package bodyio

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/http/internal/util"
)

type chunked struct {
	buf [][]byte
}

func (c *chunked) Read(p []byte) (n int, err error) {
	for len(p) > 0 && len(c.buf) > 0 {
		rn := copy(p, c.buf[0])
		n += rn

		length := len(c.buf[0])
		if length > rn {
			c.buf[0] = c.buf[0][rn:]
			return
		}
		c.buf = c.buf[1:]
	}

	if len(c.buf) == 0 {
		return 0, io.EOF
	}

	return
}

func (c *chunked) Close() error { return nil }

type ChunkedBody struct {
	buf    *bufpool.Buffer
	length int
}

func (c *ChunkedBody) Read(p []byte) (n int, err error) { panic("implement me") }
func (c *ChunkedBody) Close() error                     { panic("implement me") }

func (c *ChunkedBody) Write(p []byte) (n int, err error) {
	pLen := len(p)
	for {
		i := c.buf.Index(p)
		if i == -1 {
			_, _ = c.buf.Write(p)
			break
		}
		if i == -2 {
			c.buf.B = c.buf.B[:c.buf.Len()-1]
			p = p[1:]
		}

		_, _ = c.buf.Write(p[:i])

		if c.length < 1 {
			c.length = 1
		} else {
			c.length = -1
		}

		p = p[i+4:]
	}

	return pLen, err
}

func Chunked(preRead []byte, r io.Reader, max int) (io.ReadCloser, error) {
	br := bufio.NewReader(newBlock(preRead, r))

	bf := &chunked{}
	bf.buf = make([][]byte, max)

	var (
		n  int
		rl int
		cl [2]byte
	)

	for i := 0; ; i++ {
		line, isPrefix, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if isPrefix {
			return nil, errors.New("too long")
		}

		length := int(util.ToNonNegativeInt64(line))
		if length == 0 {
			goto END
		} else if length < 0 {
			return nil, errors.New("this is not length")
		}

		rl += length
		bf.buf[i] = make([]byte, length)

		if n, err = util.ReadAll(bf.buf[i], br); err != nil {
			return nil, err
		}
		if n != length {
			return nil, errors.New("?")
		}

	END:
		if n, err = br.Read(cl[:]); err != nil {
			return nil, err
		}
		if !bytes.Equal(cl[:], char.CRLF) {
			return nil, errors.New("what is the")
		}
		// TODO add Trailer Header
	}

	return bf, nil
}
