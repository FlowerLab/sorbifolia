package crpc

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/bufbuild/httplb"
	"github.com/bufbuild/httplb/picker"
)

var lbc = httplb.NewClient(
	httplb.WithPicker(picker.NewPowerOfTwo),
	httplb.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}, time.Second*3),
)

type gRPCClient[T any] func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) T

func Client[T any](nf gRPCClient[T], hc connect.HTTPClient, addr string, opts ...connect.ClientOption) T {
	if hc == nil {
		hc = lbc
	}

	if !strings.HasPrefix(addr, "https://") {
		addr = fmt.Sprintf("https://%s", addr)
	}

	return nf(hc, addr, append(opts, connect.WithGRPC())...)
}
