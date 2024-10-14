package flag

import (
	"reflect"
	"testing"
)

var testFlagData = []struct {
	arr []Flag
	src Flag
	dst Flag
}{
	{arr: []Flag{String | JSON | IP}, src: 0, dst: String | JSON | IP},
	{arr: []Flag{String | IP}, src: JSON, dst: String | JSON | IP},
}

func TestFlag_Set(t *testing.T) {
	for _, tt := range testFlagData {
		tt.src.Set(tt.arr...)
		if !reflect.DeepEqual(tt.src, tt.dst) {
			t.Errorf("Set(%v): got %v, want %v", tt.src, tt.dst, tt.src)
		}
	}
}
