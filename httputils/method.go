package httputils

import (
	"github.com/valyala/fasthttp"
)

func Get(uri ...string) *HTTP     { return newUtil(MethodGet, uri...) }
func Head(uri ...string) *HTTP    { return newUtil(MethodHead, uri...) }
func Post(uri ...string) *HTTP    { return newUtil(MethodPost, uri...) }
func Put(uri ...string) *HTTP     { return newUtil(MethodPut, uri...) }
func Patch(uri ...string) *HTTP   { return newUtil(MethodPatch, uri...) }
func Delete(uri ...string) *HTTP  { return newUtil(MethodDelete, uri...) }
func Options(uri ...string) *HTTP { return newUtil(MethodOptions, uri...) }
func Connect(uri ...string) *HTTP { return newUtil(MethodConnect, uri...) }
func Trace(uri ...string) *HTTP   { return newUtil(MethodTrace, uri...) }

func (h *HTTP) SetMethod(method Method) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.SetRequestURI(string(method))
		return nil
	})
}

// Method https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
type Method string

const (
	MethodGet     Method = "GET"     // RFC 7231, 4.3.1
	MethodHead    Method = "HEAD"    // RFC 7231, 4.3.2
	MethodPost    Method = "POST"    // RFC 7231, 4.3.3
	MethodPut     Method = "PUT"     // RFC 7231, 4.3.4
	MethodPatch   Method = "PATCH"   // RFC 5789
	MethodDelete  Method = "DELETE"  // RFC 7231, 4.3.5
	MethodConnect Method = "CONNECT" // RFC 7231, 4.3.6
	MethodOptions Method = "OPTIONS" // RFC 7231, 4.3.7
	MethodTrace   Method = "TRACE"   // RFC 7231, 4.3.8
)

func NewMethod(method string) Method { return Method(method) }
