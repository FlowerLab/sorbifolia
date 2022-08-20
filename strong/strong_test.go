package strong

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("", func(t *testing.T) {
		if val, err := Parse[bool]("t"); err != nil || val != true {
			t.Log("parse error")
		}
		if val, err := Parse[bool]("f"); err != nil || val != false {
			t.Log("parse error")
		}
		if _, err := Parse[bool]("ff"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[int]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[int]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[int8]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[int8]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[int16]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[int16]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[int32]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[int32]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[int64]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[int64]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[uint]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[uint]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[uint8]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[uint8]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[uint16]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[uint16]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[uint32]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[uint32]("f"); err == nil {
			t.Log("parse error")
		}
	})
	t.Run("", func(t *testing.T) {
		if val, err := Parse[uint64]("123"); err != nil || val != 123 {
			t.Log("parse error")
		}
		if _, err := Parse[uint64]("f"); err == nil {
			t.Log("parse error")
		}
	})

	defer func() { _ = recover() }()
	_ = _parse("1", map[string]string{})
	t.Error("format error")
}

func TestFormat(t *testing.T) {
	Format[bool](true)

	arr := []struct {
		t any
		v string
	}{
		{true, "true"}, {false, "false"},
		{int(1), "1"}, {uint(1), "1"},
		{int8(1), "1"}, {uint8(1), "1"},
		{int16(1), "1"}, {uint16(1), "1"},
		{int32(1), "1"}, {uint32(1), "1"},
		{int64(1), "1"}, {uint64(1), "1"},
	}
	for _, v := range arr {
		t.Run("", func(t *testing.T) {
			if _format(v.t) != v.v {
				t.Error("format error")
			}
		})
	}

	defer func() { _ = recover() }()
	_ = _format(map[struct{}]string{})
	t.Error("format error")
}

func TestAppend(t *testing.T) {
	Append(nil, true)
	arr := []struct {
		t any
		v string
	}{
		{true, "true"}, {false, "false"},
		{int(1), "1"}, {uint(1), "1"},
		{int8(1), "1"}, {uint8(1), "1"},
		{int16(1), "1"}, {uint16(1), "1"},
		{int32(1), "1"}, {uint32(1), "1"},
		{int64(1), "1"}, {uint64(1), "1"},
	}
	for _, v := range arr {
		t.Run("", func(t *testing.T) {
			if !bytes.Equal(_append(nil, v.t), []byte(v.v)) {
				t.Error("append error")
			}
		})
	}

	defer func() { _ = recover() }()
	_ = _append(nil, map[struct{}]string{})
	t.Error("append error")
}
