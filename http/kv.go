package http

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
)

type KV struct {
	K    []byte
	V    []byte
	Null bool
}

var nullKV = KV{Null: true}

func (kv *KV) SetK(b []byte) { kv.K = append(kv.K, b...) }
func (kv *KV) SetV(b []byte) { kv.V = append(kv.V, b...) }

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
			kv.V = b[idx:]
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

func (ks *KVs) Find(key []byte) []KV {
	arr := make([]int, 0, len(*ks))
	for i, v := range *ks {
		if bytes.EqualFold(key, v.K) {
			arr = append(arr, i)
		}
	}
	if len(arr) == 0 {
		return nil
	}

	kvs := make([]KV, 0, len(arr))
	for i := range arr {
		kvs = append(kvs, (*ks)[i])
	}
	return kvs
}

func (ks *KVs) add(kv KV) {
	*ks = append(*ks, kv)
}

func (ks *KVs) preAlloc(size int) {
	var l = len(*ks)
	if size <= 0 {
		size = 1
	}

	if cap(*ks) < l+size {
		*ks = append(*ks, make([]KV, size)...)
		*ks = (*ks)[:l]
	}
}

func (ks *KVs) alloc() *KV {
	var l = len(*ks)
	if cap(*ks) > l {
		*ks = (*ks)[:l+1]
	} else {
		*ks = append(*ks, KV{})
	}
	return &(*ks)[l]
}

func (ks *KVs) addHeader(b []byte) {
	kv := ks.alloc()
	idx := bytes.IndexByte(b, char.Colon)
	if idx == -1 {
		kv.SetK(b)
		return
	}

	kv.K = append(kv.K, b[:idx]...)
	kv.SetK(b[:idx])
	idx++
	for ; idx < len(b); idx++ {
		if b[idx] != char.Space {
			kv.SetV(b[idx:])
			break
		}
	}
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
