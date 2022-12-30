package http

import (
	"arena"
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
)

type KV struct {
	K []byte
	V *[]byte
}

var nullKV = KV{}

func (kv KV) IsNil() bool         { return len(kv.K) == 0 }
func (kv KV) Equal(b []byte) bool { return bytes.Equal(kv.Val(), b) }

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
	idx++
	for ; idx < len(b); idx++ {
		if b[idx] != char.Space {
			v := b[idx:]
			kv.V = &v
			break
		}
	}
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

func (ks *KVs) Find(key []byte, a *arena.Arena) []KV {
	arr := arena.MakeSlice[int](a, 0, len(*ks))
	for i, v := range *ks {
		if bytes.EqualFold(key, v.K) {
			arr = append(arr, i)
		}
	}
	if len(arr) == 0 {
		return nil
	}

	kvs := arena.MakeSlice[KV](a, 0, len(arr))
	for i := range arr {
		kvs = append(kvs, (*ks)[i])
	}
	return kvs
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
