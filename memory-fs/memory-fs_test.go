package mfs

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestPersistence(t *testing.T) {
	fs := New()
	err := fs.(*mfs).MkdirAll("a/b/c")
	if err != nil {
		t.Error(err)
	}

	err = Persistence(fs, "D:\\pic")
	if err != nil {
		t.Error(err)
	}
}

func TestFork(t *testing.T) {
	fs, err := Fork("D:/pic/a/b/c")
	if err != nil {
		t.Error(err)
	}

	mf := fs.(*mfs)
	d, ok := mf.root.node["pic"]
	if !ok {
		t.Error("fail to fork")
	}

	f := &file{
		name:    "test.txt",
		perm:    0,
		modTime: time.Now(),
		data:    []byte("123456dd"),
	}
	d.(*dir).node["test.txt"] = f
	if err = Persistence(fs, "D:"); err != nil {
		t.Error(err)
	}

	tt, err := os.Open("D:/pic/test.txt")
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(tt); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(buf.Bytes(), f.data) {
		t.Errorf("expect: %v,get: %v", buf.Bytes(), f.data)
	}
}
