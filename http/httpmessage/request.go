package httpmessage

import (
	"io"

	"go.x2ox.com/sorbifolia/http/httpconfig"
	"go.x2ox.com/sorbifolia/http/httpheader"
	"go.x2ox.com/sorbifolia/http/internal/bufpool"
	"go.x2ox.com/sorbifolia/http/method"
	"go.x2ox.com/sorbifolia/http/version"
)

type Request struct {
	cfg        *httpconfig.Config
	state      state
	buf        *bufpool.ReadBuffer
	p          int
	bodyLength int

	Version version.Version
	Method  method.Method
	Header  httpheader.RequestHeader
	Body    io.ReadWriteCloser
}
