package mfs

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestPersistence(t *testing.T) {
	var (
		home      = os.Getenv("HOME")
		targetDir = home + "/pic"
	)

	if !exists(targetDir) {
		if err := os.Mkdir(targetDir, 0777); err != nil {
			t.Error(err)
		}
	}

	fs := New()
	err := fs.(*mfs).MkdirAll("a/b/c")
	if err != nil {
		t.Error(err)
	}

	mf := fs.(*mfs)
	mf.root.node["test.txt"] = &file{
		name:    "test.txt",
		perm:    os.ModePerm,
		modTime: time.Now(),
		data:    []byte("1253test"),
	}

	err = Persistence(fs, targetDir)
	if err != nil {
		t.Error(err)
	}

	var paths = []string{
		targetDir + "/a",
		targetDir + "/a/b",
		targetDir + "/a/b/c",
		targetDir + "/test.txt",
	}

	for _, p := range paths {
		if !exists(p) {
			t.Error("fail to persistence")
		}
	}
}

func TestFork(t *testing.T) {
	test := []struct {
		dp        string
		df        string
		targetDir string
		m         *mfs
		mf        *dir
		Err       error
	}{
		{
			"pic/a/b/c",
			"pic/ttt.txt",
			"/pic",
			&mfs{},
			&dir{},
			nil,
		},
	}

	for i, v := range test {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if v.targetDir[0] == '/' {
				v.targetDir = v.targetDir[1:]
			}

			if !exists(v.dp) {
				if err := os.MkdirAll(v.targetDir, 0777); err != nil {
					t.Error(err)
				}
			}
			if !exists(v.df) {
				if _, err := os.Create(v.df); err != nil {
					t.Error(err)
				}
			}

			_, err := Fork(v.targetDir)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}

			_ = os.RemoveAll(v.targetDir)
			_ = os.RemoveAll(v.df)
		})
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
