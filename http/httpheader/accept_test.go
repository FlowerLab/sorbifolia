package httpheader

import (
	"fmt"
	"testing"
)

func TestContentLength_Length(t *testing.T) {
	v := ContentLength("12300800")
	fmt.Println(v.Length())
}

func BenchmarkContentLength_Length(b *testing.B) {
	b.ReportAllocs()

	v := ContentLength("12300800")
	for i := 0; i < b.N; i++ {
		v.Length()
	}
}
