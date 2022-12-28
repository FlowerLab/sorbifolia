package httpheader

import (
	"bytes"
	"strconv"

	"go.x2ox.com/sorbifolia/http/internal/char"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

func cleanSuffixSpace(b []byte) []byte {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != char.Space {
			return b[:i]
		}
	}
	return nil
}

func cleanPrefixSpace(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		if b[i] != char.Space {
			return b[i:]
		}
	}
	return nil
}

func cleanQuotationMark(b []byte) []byte {
	if len(b) < 2 || b[0] != char.QuotationMark || b[len(b)-1] != char.QuotationMark {
		return b
	}
	return b[1 : len(b)-1]
}

func eachQualityValue(b []byte, fn EachQualityValue) {
	eachValueWithComma(b, func(val []byte) bool {
		qv := QualityValue{Value: val, Priority: 1}

		if i := bytes.IndexByte(val, char.Semi); i >= 0 {
			qv.Value = val[:i]

			eachValueWithSemi(val[i+1:], func(value []byte) bool {
				if k, v := parseKVWithEqual(value); len(k) == 1 && k[0] == 'q' {
					qv.Priority, _ = strconv.ParseFloat(pyrokinesis.Bytes.ToString(v), 64)
					return false
				}
				return true
			})
		}

		return fn(qv)
	})
}

func eachValueWithComma(b []byte, fn EachValue) { eachValue(b, char.Comma, true, fn) }
func eachValueWithSemi(b []byte, fn EachValue)  { eachValue(b, char.Semi, true, fn) }

func parseKVWithEqual(b []byte) (key, val []byte) {
	if i := bytes.IndexByte(b, char.Equal); i >= 0 {
		return b[:i], b[i+1:]
	}
	return b, nil
}

func toNonNegativeInt64(b []byte) (n int64) {
	for _, val := range b {
		if val > '9' || val < '0' {
			return 0
		}
		n = n*10 + int64(val-'0')
	}
	return
}

func getHeaderValue(b []byte) []byte { return parseFirstValue(b, char.Semi) }
func getHeaderParam(b []byte) []byte { return parseOtherValue(b, char.Semi) }

func parseFirstValue(b []byte, delimiter byte) []byte {
	if i := bytes.IndexByte(b, delimiter); i >= 0 {
		return b[:i]
	}
	return b
}
func parseOtherValue(b []byte, delimiter byte) []byte {
	if i := bytes.IndexByte(b, delimiter); i >= 0 {
		return b[i+1:]
	}
	return b
}

func eachValue(b []byte, delimiter byte, cleanSpace bool, fn EachValue) {
	for {
		if len(b) == 0 {
			return
		}

		i := bytes.IndexByte(b, delimiter)
		if i < 0 {
			i = len(b)
		}

		val := b[:i]

		if len(b) < i+1 {
			i = len(b) - 1
		}
		b = b[i+1:]

		if cleanSpace {
			val = cleanPrefixSpace(val)
		}

		if !fn(val) {
			return
		}
	}
}
