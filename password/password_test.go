package password

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync/atomic"
	"testing"
)

var (
	testNum = int32(4)
	testCh  atomic.Int32
)

func TestMustGenerate(t *testing.T) {
	t.Parallel()

	hash := MustGenerate("password")
	if !Compare(hash, "password") {
		t.Errorf("%s | %s not match", "password", hash)
	}
	testCh.Add(1)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	hash, _ := Generate("password")
	if !Compare(hash, "password") {
		t.Errorf("%s | %s not match", "password", hash)
	}
	testCh.Add(1)
}

func TestCompare(t *testing.T) {
	t.Parallel()

	hash := MustGenerate("password")
	if !Compare(hash, "password") {
		t.Errorf("%s | %s not match", "password", hash)
	}
	testCh.Add(1)
}

func TestFail(t *testing.T) {
	t.Parallel()

	t.Run("", func(t *testing.T) {
		if Compare("", "") {
			t.Error("fail")
		}
	})

	t.Run("", func(t *testing.T) {
		if Compare(base64.RawStdEncoding.EncodeToString([]byte("1234567890")), "") {
			t.Error("fail")
		}
	})

	t.Run("", func(t *testing.T) {
		if Compare("AAAAAQABAAABAAAALZP/dD6HbO0SPK8Zijd/ivOT/3G3Wj1SKzrkIKs3REnw", "1") {
			t.Error("fail")
		}
	})

	testCh.Add(1)
}

type errReader struct{}

func (e errReader) Read([]byte) (n int, err error) {
	return 0, errors.New("OEF")
}

func TestRandReaderErr(t *testing.T) {
	t.Parallel()

	for {
		if testCh.CompareAndSwap(testNum, -1) {
			break
		}
	}

	g := New()

	rand.Reader = errReader{}
	defer func() { _ = recover() }()

	g.MustGenerate("123456")
	t.Error("fail")
}
