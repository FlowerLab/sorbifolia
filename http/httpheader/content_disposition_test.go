package httpheader

import (
	"fmt"
	"testing"
)

func TestContentDisposition_Param(t *testing.T) {
	cd := ContentDisposition("attachment; filename=\"name.ext\"; name=123")
	fmt.Println(string(cd.Param([]byte("filename"))))
}

func BenchmarkA(b *testing.B) {
	b.ReportAllocs()

	cd := ContentDisposition("attachment; filename=\"name.ext\"; name=123")
	k := []byte("filename")
	for i := 0; i < b.N; i++ {
		cd.Param(k)
	}
}
