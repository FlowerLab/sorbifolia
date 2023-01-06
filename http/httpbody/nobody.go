package httpbody

import (
	"io"
)

var _nobody = nobody{}

type nobody struct{}

func (nobody) release()                    {}
func (nobody) Reset()                      {}
func (nobody) Write(_ []byte) (int, error) { return 0, io.EOF }
func (nobody) Read(_ []byte) (int, error)  { return 0, io.EOF }
func (nobody) Close() error                { return nil }

func Null() io.ReadWriteCloser { return _nobody }
