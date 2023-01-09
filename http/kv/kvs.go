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

func (ks *KVs) HasKey(key []byte) bool { return ks.Index(key) != -1 }

func (ks *KVs) Get(key []byte) *KV {
	if i := ks.Index(key); i != -1 {
		return &(*ks)[i]
	}
	return nil
}
func (ks *KVs) GetValue(key []byte) []byte {
	if i := ks.Index(key); i != -1 {
		return (*ks)[i].V
	}
	return nil
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

func (ks *KVs) Add(k, v []byte) { kv := ks.alloc(); kv.K, kv.V = k, v }
func (ks *KVs) AddKV(kv KV)     { *ks = append(*ks, kv) }
func (ks *KVs) Set(k, v []byte) {
	if i := ks.Index(k); i != -1 {
		(*ks)[i].V, (*ks)[i].Null = v, false
		return
	}
	ks.Add(k, v)
}
func (ks *KVs) SetKV(kv KV) {
	if i := ks.Index(kv.K); i != -1 {
		(*ks)[i].V, (*ks)[i].Null = kv.V, kv.Null
		return
	}
	ks.AddKV(kv)
}

func (ks *KVs) Index(k []byte) int {
	for i := range *ks {
		if bytes.EqualFold((*ks)[i].K, k) {
			return i
		}
	}
	return -1
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
