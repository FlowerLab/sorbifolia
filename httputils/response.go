package httputils

import (
	"github.com/valyala/fasthttp"
)

func (h *HTTP) ParserData(arr ...Parser) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		for _, handle := range arr {
			if err := handle(resp); err != nil {
				return err
			}
		}
		return nil
	})
}
