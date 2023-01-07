package kv

import (
	"bytes"

	"go.x2ox.com/sorbifolia/http/internal/char"
)

type KVs []KV

func (ks *KVs) Reset() {
	for i := range *ks {
		(*ks)[i].Reset()
	}
	*ks = (*ks)[:0]
}

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

func (ks *KVs) Add(kv KV) {
	*ks = append(*ks, kv)
}

func (ks *KVs) PreAlloc(size int) {
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

func (ks *KVs) AddHeader(b []byte) {
	kv := ks.alloc()
	idx := bytes.IndexByte(b, char.Colon)
	if idx == -1 {
		kv.SetK(b)
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

func (ks *KVs) Set(kv KV) {
	for _, v := range *ks {
		if bytes.EqualFold(v.K, kv.K) {
			v.V = kv.V
			return
		}
	}
	ks.Add(kv)
}

var nullKV = KV{Null: true}
