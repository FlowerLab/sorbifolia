package method

type Method string

const (
	Get     Method = "GET"
	Head    Method = "HEAD"
	Post    Method = "POST"
	Put     Method = "PUT"
	Patch   Method = "PATCH" // RFC 5789
	Delete  Method = "DELETE"
	Connect Method = "CONNECT"
	Options Method = "OPTIONS"
	Trace   Method = "TRACE"
)
