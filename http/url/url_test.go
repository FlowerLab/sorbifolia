package url

import (
	"fmt"
	"testing"
)

func TestURL_Bytes(t *testing.T) {
	u := &URL{}
	if err := u.Parse([]byte("google.com:123"), []byte("/123/asd?asa=aa#zxc"), true); err != nil {
		t.Error(err)
	}
	// fmt.Printf("%s", string(u.Scheme))
	// u.full = append(u.full[:0],
	// 	u
	// 	)
	fmt.Println(string(u.Bytes()))
	fmt.Println(string(u.FullBytes()))
}
