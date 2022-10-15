package strutils

import (
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

func Append(dst, src string) string {
	db := pyrokinesis.String.ToBytes(dst)
	return pyrokinesis.Bytes.ToString(AppendToBytes(db, src))
}

func AppendToBytes(dst []byte, src string) []byte {
	sb := pyrokinesis.String.ToBytes(src)

	if l := len(dst) + len(sb); cap(dst) < l {
		out := make([]byte, l)
		copy(out, dst)
		copy(out[len(dst):], sb)
		return out
	}
	return append(dst, sb...)
}
