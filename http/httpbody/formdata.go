package httpbody

import (
	"go.x2ox.com/sorbifolia/http/kv"
)

// FormData multipart/form-data
type FormData struct {
	Boundary []byte

	KV   kv.KVs
	File []FileHeader
}

type FileHeader struct {
	Name     []byte
	Filename []byte
	Size     int64
	Header   kv.KVs

	content []byte
	tmpfile string
}
