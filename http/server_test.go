package http

import (
	"net"
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
			ctx.Response.StatusCode = 200
			ctx.Response.SetBody("asdsdaasdas")
		},
	}
	// go http.ListenAndServe("127.0.0.1:6060", nil)

	ln, _ := net.Listen("tcp", "127.0.0.1:8808")
	s.Serve(ln)
}
