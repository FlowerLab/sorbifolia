package padding

import (
	"bytes"
	"testing"
)

func TestPKCS7_Pad(t *testing.T)         { testPKCS7(t) }
func TestPKCS7_UnPad(t *testing.T)       { testPKCS7(t) }
func TestNoPadding_Pad(t *testing.T)     { testNoPadding(t) }
func TestNoPadding_UnPad(t *testing.T)   { testNoPadding(t) }
func TestZeroPadding_Pad(t *testing.T)   { testZeroPadding(t) }
func TestZeroPadding_UnPad(t *testing.T) { testZeroPadding(t) }
func TestISO10126_Pad(t *testing.T)      { testISO10126(t) }
func TestISO10126_UnPad(t *testing.T)    { testISO10126(t) }
func TestANSIx923_Pad(t *testing.T)      { testANSIx923(t) }
func TestANSIx923_UnPad(t *testing.T)    { testANSIx923(t) }

func testPKCS7(t *testing.T) {
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
}

func testNoPadding(t *testing.T) {
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
}

func testZeroPadding(t *testing.T) {
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
}

func testISO10126(t *testing.T) {
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
}
func testANSIx923(t *testing.T) {
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
}
