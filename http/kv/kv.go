package kv

type KV struct {
	K, V []byte
	Null bool
}

func (kv *KV) SetK(b []byte) { kv.K = append(kv.K[:0], b...) }
func (kv *KV) SetV(b []byte) { kv.V = append(kv.V[:0], b...) }
func (kv *KV) Reset()        { kv.K = kv.K[:0]; kv.V = kv.V[:0]; kv.Null = false }
