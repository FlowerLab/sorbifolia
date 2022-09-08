package pyrokinesis

import (
	"strconv"
	"testing"
)

func TestBytes_ToString(t *testing.T) {
	Bytes.ToString([]byte("hello"))
}

func Test_Bytes_ToNumber(t *testing.T) {
	if strconv.IntSize == 64 {
		if n := Bytes.ToInt([]byte{123, 123, 12, 0, 0, 0, 0, 255}); n != -72057594037109893 {
			t.Error("Bytes.ToNumber err num:", n)
		}
		if n := Bytes.ToUint([]byte{123, 123, 12, 0, 0, 0, 0, 255}); n != 18374686479672441723 {
			t.Error("Bytes.ToNumber err num:", n)
		}
	}
	if n := Bytes.ToInt8([]byte{123}); n != 123 {
		t.Error("Bytes.ToNumber err num:", n)
	}
	if n := Bytes.ToInt16([]byte{123, 255}); n != -133 {
		t.Error("Bytes.ToNumber err num:", n)
	}
	if n := Bytes.ToInt32([]byte{123, 123, 123, 255}); n != -8684677 {
		t.Error("Bytes.ToNumber err num:", n)
	}
	if n := Bytes.ToInt64([]byte{123, 123, 12, 0, 0, 0, 0, 255}); n != -72057594037109893 {
		t.Error("Bytes.ToNumber err num:", n)
	}

	if n := Bytes.ToUint8([]byte{123}); n != 123 {
		t.Error("Bytes.ToNumber err num:", n)
	}
	if n := Bytes.ToUint16([]byte{123, 255}); n != 65403 {
		t.Error("Bytes.ToNumber err num:", n)
	}
	if n := Bytes.ToUint32([]byte{123, 123, 123, 255}); n != 4286282619 {
		t.Error("Bytes.ToNumber err num:", n)
	}
	if n := Bytes.ToUint64([]byte{123, 123, 12, 0, 0, 0, 0, 255}); n != 18374686479672441723 {
		t.Error("Bytes.ToNumber err num:", n)
	}
}
