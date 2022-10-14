package coarsetime

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

var (
	testStr1a = "阿ab12三啊11232123实打123实121233312打3算123"
	testStr1b = "321算3打213332121实321打实32123211啊三21ba阿"
)

func TestReverse(t *testing.T) {
	if Reverse(testStr1a) != testStr1b {
		t.Error("fail")
	}
}

func TestReverseB(t *testing.T) {
	if ReverseB(testStr1a) != testStr1b {
		t.Error("fail")
	}
}

func TestSAD(t *testing.T) {
	b := []byte("a")
	r, s := utf8.DecodeRune(b)
	fmt.Printf("%08b, %s, %d, %d\n", b, string(r), s, checkUTF8Len(b))

	b = []byte("é")
	r, s = utf8.DecodeRune(b)
	fmt.Printf("%08b, %s, %d, %d\n", b, string(r), s, checkUTF8Len(b))

	b = []byte("我")
	r, s = utf8.DecodeRune(b)
	fmt.Printf("%08b, %s, %d, %d\n", b, string(r), s, checkUTF8Len(b))

	b = []byte(string('\U0010FFFF'))
	r, s = utf8.DecodeRune(b)
	fmt.Printf("%08b, %s, %d, %d\n", b, string(r), s, checkUTF8Len(b))

	bbb := []byte(testStr1a)
	fmt.Println(string(bbb[:3]), len(bbb[:4]))
	checkLastUTF8Len(bbb[:4])

}

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reverse(testStr1a)
	}
}

func BenchmarkReverseB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReverseB(testStr1a)
	}
}
