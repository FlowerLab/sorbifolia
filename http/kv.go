package http

import (
	"bytes"
	"strconv"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/pyrokinesis"
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
	idx++
	for ; idx < len(b); idx++ {
		if b[idx] != char.Space {
			v := b[idx:]
			kv.V = &v
			break
		}
	}
}

type QualityValue struct {
	Value    []byte
	Priority float64 // 1.00 - 0.00
}

func (kv KV) QualityValues(b []byte) QualityValue {
	if kv.V == nil {
		return QualityValue{Priority: -1}
	}
	buf := *kv.V

	for {
		if len(buf) == 0 {
			return QualityValue{Priority: -1}
		}

		i := bytes.IndexByte(buf, char.Comma)
		if i < 0 {
			i = len(buf)
		}

		val := buf[:i]

		{
			if len(buf) < i+1 {
				i = len(buf) - 1
			}
			buf = buf[i+1:]
		}

		{
			j := 0
			for ; j < len(val); j++ {
				if val[j] != char.Space {
					break
				}
			}
			val = val[j:]
		}

		var qv QualityValue
		if i = bytes.IndexByte(val, char.Semi); i < 0 {
			qv.Value = val
			qv.Priority = 1
		} else {
			qv.Value = val[:i]
			val = val[i:]
			if i = bytes.IndexByte(val, '='); i > 0 {
				qv.Priority, _ = strconv.ParseFloat(pyrokinesis.Bytes.ToString(val[i+1:]), 64)
			}
		}

		if bytes.EqualFold(qv.Value, b) {
			return qv
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

// Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5
// Set-Cookie: UserID=JohnDoe; Max-Age=3600; Version=1
// Content-Type: text/html; charset=utf-8
// Content-Disposition: attachment; filename="name.ext"
