package httputils

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type Parser func(resp *fasthttp.Response) error

func JSONParser(v any) Parser {
	return func(resp *fasthttp.Response) error {
		data, err := resp.BodyUncompressed()
		if err != nil {
			return err
		}
		return json.Unmarshal(data, v)
	}
}

func HeaderParser(m map[string]string) Parser {
	return func(resp *fasthttp.Response) error {
		for k := range m {
			m[k] = string(resp.Header.Peek(k))
		}
		return nil
	}
}
