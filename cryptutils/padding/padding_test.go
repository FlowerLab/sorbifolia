package padding

import (
	"bytes"
	"crypto/rand"
	"errors"
	"testing"
)

func TestPKCS7(t *testing.T) {
	var p PKCS7
	data := []byte{1, 2, 3, 4, 5}

	padded, err := p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(padded, []byte{1, 2, 3, 4, 5, 3, 3, 3}) {
		t.Fatalf("Wrong padding")
	}

	data = []byte{1, 2, 3, 4, 5, 6, 7, 8}

	padded, err = p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(padded, []byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 8}) {
		t.Fatalf("Wrong padding")
	}

	if _, err = p.UnPad(padded, 8); err != nil {
		t.Error("fail")
	}
	if _, err = p.UnPad(padded, 9); err == nil {
		t.Error("fail")
	}
	if _, err = p.UnPad([]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 7}, 8); err == nil {
		t.Error("fail")
	}
	if _, err = p.UnPad([]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 9}, 8); err == nil {
		t.Error("fail")
	}
	if _, err = p.UnPad([]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 1, 8}, 8); err == nil {
		t.Error("fail")
	}
}

func TestNoPadding(t *testing.T) {
	var p NoPadding
	data := []byte{1, 2, 3, 4, 5}

	padded, err := p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(padded, data) {
		t.Fatalf("Wrong padding")
	}

	data = []byte{1, 2, 3, 4, 5, 6, 7, 8}

	padded, err = p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(padded, data) {
		t.Fatalf("Wrong padding")
	}

	if _, err = p.UnPad(padded, 8); err != nil {
		t.Error("fail")
	}
}

func TestZeroPadding(t *testing.T) {
	var p ZeroPadding
	data := []byte{1, 2, 3, 4, 5}

	padded, err := p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(padded, []byte{1, 2, 3, 4, 5, 0, 0, 0}) {
		t.Fatalf("Wrong padding")
	}

	data = []byte{1, 2, 3, 4, 5, 6, 7, 8}

	padded, err = p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(padded, []byte{1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0}) {
		t.Fatalf("Wrong padding")
	}

	if _, err = p.UnPad(padded, 8); err != nil {
		t.Error("fail")
	}
	if _, err = p.UnPad(padded, 9); err == nil {
		t.Error("fail")
	}
	if _, err = p.UnPad([]byte{1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 8); err == nil {
		t.Error("fail")
	}
}

func TestISO10126(t *testing.T) {
	var p ISO10126
	data := []byte{1, 2, 3, 4, 5}

	padded, err := p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if padded[7] != byte(3) {
		t.Fatalf("Wrong padding")
	}

	data = []byte{1, 2, 3, 4, 5, 6, 7, 8}

	padded, err = p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if padded[15] != byte(8) {
		t.Fatalf("Wrong padding")
	}

	if _, err = p.UnPad(padded, 8); err != nil {
		t.Error("fail")
	}
	if _, err = p.UnPad(padded, 9); err == nil {
		t.Error("fail")
	}
	if _, err = p.UnPad([]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 9}, 8); err == nil {
		t.Error("fail")
	}

	rand.Reader = errReader{}
	if _, err = p.Pad(data, 8); err == nil {
		t.Fatal("err")
	}
}

func TestANSIx923(t *testing.T) {
	var p ANSIx923
	data := []byte{1, 2, 3, 4, 5}

	padded, err := p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(padded, []byte{1, 2, 3, 4, 5, 0, 0, 3}) {
		t.Fatalf("Wrong padding")
	}

	data = []byte{1, 2, 3, 4, 5, 6, 7, 8}

	padded, err = p.Pad(data, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(padded, []byte{1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 8}) {
		t.Fatalf("Wrong padding")
	}

	if _, err = p.UnPad(padded, 8); err != nil {
		t.Error("fail")
	}
	if _, err = p.UnPad(padded, 9); err == nil {
		t.Error("fail")
	}
	if _, err = p.UnPad([]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 7}, 8); err == nil {
		t.Error("fail")
	}
	if _, err = p.UnPad([]byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 9}, 8); err == nil {
		t.Error("fail")
	}
}

func TestPadding(t *testing.T) {
	t.Run("", func(t *testing.T) {
		for _, v := range []Padding{PKCS7{}, ZeroPadding{}, ISO10126{}, ANSIx923{}} {
			if _, err := v.Pad(nil, -1); err == nil {
				t.Error("err")
			}
		}
	})
}

type errReader struct{}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("OEF")
}
