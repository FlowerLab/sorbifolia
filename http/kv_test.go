//go:build goexperiment.arenas

package http

import (
	"testing"
)

func TestKV_QualityValues(t *testing.T) {
	testData := []struct {
		s string
		t []QualityValue
	}{
		{
			s: "Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
			t: []QualityValue{
				{[]byte("fr-CH"), 1},
				{[]byte("fr"), 0.9},
			},
		},
	}

	for _, v := range testData {
		t.Run("", func(t *testing.T) {
			var kv KV
			kv.ParseHeader([]byte(v.s))

			for _, vv := range v.t {

				if kv.QualityValues(vv.Value).Priority != vv.Priority {
					qv := kv.QualityValues(vv.Value)
					t.Errorf("asd %v %s %s", qv, string(qv.Value), vv.Value)
				}
			}
		})
	}
}
