package method

import (
	"bytes"

	"go.x2ox.com/sorbifolia/pyrokinesis"
)

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

var methods = map[string]Method{
	"GET": Get, "HEAD": Head, "POST": Post, "PUT": Put, "PATCH": Patch,
	"DELETE": Delete, "CONNECT": Connect, "OPTIONS": Options, "TRACE": Trace,
}

func Parse(b []byte) Method {
	s := pyrokinesis.Bytes.ToString(b)
	if m, ok := methods[s]; ok {
		return m
	}
	return Method(s)
}

func (m Method) IsGet() bool     { return m.Is(Get) }
func (m Method) IsHead() bool    { return m.Is(Head) }
func (m Method) IsPost() bool    { return m.Is(Post) }
func (m Method) IsPut() bool     { return m.Is(Put) }
func (m Method) IsPatch() bool   { return m.Is(Patch) }
func (m Method) IsDelete() bool  { return m.Is(Delete) }
func (m Method) IsConnect() bool { return m.Is(Connect) }
func (m Method) IsOptions() bool { return m.Is(Options) }
func (m Method) IsTrace() bool   { return m.Is(Trace) }
func (m Method) Is(method Method) bool {
	return bytes.EqualFold(
		pyrokinesis.String.ToBytes(string(m)),
		pyrokinesis.String.ToBytes(string(method)),
	)
}

func (m Method) Bytes() []byte {
	if len(m) == 0 {
		return pyrokinesis.String.ToBytes(string(Get))
	}
	return pyrokinesis.String.ToBytes(string(m))
}
