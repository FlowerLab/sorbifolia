package kv

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
)

type KV struct {
	K    []byte
	V    []byte
	Null bool
}

func (kv *KV) SetK(b []byte) { kv.K = append(kv.K, b...) }
func (kv *KV) SetV(b []byte) { kv.V = append(kv.V, b...) }

func (kv *KV) Reset() { kv.K = kv.K[:0]; kv.V = kv.V[:0]; kv.Null = false }

func (kv *KV) ParseHeader(b []byte) {
	idx := bytes.IndexByte(b, char.Colon)
	if idx == -1 {
		kv.SetK(b)
		kv.Null = true
		return
	}

	kv.SetK(b[:idx])
	idx++
	for ; idx < len(b); idx++ {
		if b[idx] != char.Space {
			kv.SetV(b[idx:])
			break
		}
	}
}

func (kv *KV) AppendHeader(dst []byte) []byte {
	if dst = append(dst, kv.K...); kv.Null {
		return dst
	}
	dst = append(dst, char.Colon)
	dst = append(dst, char.Space)
	dst = append(dst, kv.V...)
	dst = append(dst, char.CRLF...)
	return dst
}
