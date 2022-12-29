package bodyio

import (
	"arena"
	"bytes"
	"errors"
	"io"

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

func Chunked(a *arena.Arena, preRead []byte, r io.Reader, max int) (io.ReadCloser, error) {
	br := arena.New[util.Reader](a)
	br.Reset(newBlock(a, preRead, r), arena.MakeSlice[byte](a, 1024, 1024))

	bf := arena.New[chunked](a)
	bf.buf = arena.MakeSlice[[]byte](a, max, max)

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
		bf.buf[i] = arena.MakeSlice[byte](a, length, length)

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
	}

	return bf, nil
}
