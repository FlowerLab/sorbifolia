package method

import (
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
