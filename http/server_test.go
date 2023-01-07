package http

import (
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"

	"go.x2ox.com/sorbifolia/http/httpconfig"
)

func TestS(t *testing.T) {
	s := &Server{
		Config: httpconfig.Config{
			MaxIdleWorkerDuration: time.Second * 30,
		},

		Handler: func(ctx *Context) {
			var b []byte
			if ctx.Request.Body != nil {
				if b, _ = io.ReadAll(ctx.Request.Body); len(b) == 0 {
					b = []byte("nobody nobody")
				}
			}
			ctx.Response.StatusCode = 200
			ctx.Response.SetBody(b)
		},
	}
	go http.ListenAndServe("127.0.0.1:6060", nil)

	ln, _ := net.Listen("tcp", "127.0.0.1:8808")
	s.Serve(ln)
}
