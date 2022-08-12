package httputils

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type ContentType string

func (c ContentType) SetCharset(charset string) ContentType {
	return NewContentType(fmt.Sprintf("%s; charset=%s", c, charset))
}

func (c ContentType) SetBoundary(boundary string) ContentType {
	return NewContentType(fmt.Sprintf("%s; boundary=%s", c, boundary))
}

func NewContentType(contentType string) ContentType { return ContentType(contentType) }

const (
	TextHTML          ContentType = "text/html"
	TextPlain         ContentType = "text/plain"
	AppOctetStream    ContentType = "application/octet-stream"
	AppJSON           ContentType = "application/json"
	AppFormUrlencoded ContentType = "application/x-www-form-urlencoded"
	AppXML            ContentType = "application/xml"
	MultiFormData     ContentType = "multipart/form-data"
)

func (h *HTTP) SetContentType(contentType ContentType) *HTTP {
	return h.Add(func(client *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
		req.Header.SetContentType(string(contentType))
		return nil
	})
}
