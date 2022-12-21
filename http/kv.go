package http

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
)

type KV struct {
	K []byte
	V *[]byte
}

var nullKV = KV{}

func (kv KV) IsNil() bool { return len(kv.K) == 0 }

func (kv KV) Val() []byte {
	if kv.V == nil {
		return nil
	}
	return *kv.V
}

func (kv *KV) ParseHeader(b []byte) {
	idx := bytes.IndexByte(b, char.Colon)
	if idx == -1 {
		kv.K = b
		return
	}

	kv.K = b[:idx]

	length := len(b)

	for {
		idx++
		if length < idx {
			break
		}
		if b[idx] != char.Space {
			v := b[idx:]
			kv.V = &v
			break
		}
	}
}

type QualityValues struct {
	Value    []byte
	Priority float32 // 1.00 - 0.00
}

func (kv *KV) List(b []byte) {

}

type KVs []KV

func (ks *KVs) Each(fn func(kv KV) bool) {
	for _, v := range *ks {
		if !fn(v) {
			return
		}
	}
}

func (ks *KVs) HasKey(key []byte) bool {
	for _, v := range *ks {
		if bytes.EqualFold(key, v.K) {
			return true
		}
	}
	return false
}

func (ks *KVs) Get(key []byte) KV {
	for _, v := range *ks {
		if bytes.EqualFold(key, v.K) {
			return v
		}
	}
	return nullKV
}

func (ks *KVs) add(kv KV) {
	*ks = append(*ks, kv)
}

func (ks *KVs) set(kv KV) {
	for _, v := range *ks {
		if bytes.EqualFold(v.K, kv.K) {
			v.V = kv.V
			return
		}
	}
	ks.add(kv)
}

// zh-CN,
// zh;q=0.8,
// zh-TW;q=0.7,
// zh-HK;q=0.5,
// en-US;q=0.3,
// en;q=0.2
