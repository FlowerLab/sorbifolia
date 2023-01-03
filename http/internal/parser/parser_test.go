package parser

import (
	"fmt"
	"testing"
)

func TestRequestParser_parseHeader(t *testing.T) {
	rp := &RequestParser{SetHeaders: func(b []byte) error {
		fmt.Println(string(b))
		return nil
	}}

	rp.parseHeader([]byte("A:b\r"))
	rp.parseHeader([]byte("\n\r\n"))

}
