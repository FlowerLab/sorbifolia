package mfs

import (
	"fmt"
	"os"
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

	tests := []struct {
		td string
		MemoryFS
		paths []string
		Err   error
	}{
		{
			"/pic/a",
			fs,
			nil,
			&os.PathError{},
		},
		{
			targetDir,
			fs,
			[]string{
				targetDir + "/a",
				targetDir + "/a/b",
				targetDir + "/a/b/c",
				targetDir + "/test.txt",
			},
			nil,
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err = Persistence(v.MemoryFS, v.td)
			if err != nil && v.Err == nil {
				t.Error(err)
			}
			if err == nil {
				for _, p := range v.paths {
					if !exists(p) {
						t.Error("fail to persistence")
					}
				}
				_ = os.RemoveAll(v.td)
			}
		})
	}
}

func TestFork(t *testing.T) {
	test := []struct {
		dp        string
		df        string
		targetDir string
		rmf       string
		m         *mfs
		mf        *dir
		Err       error
	}{
		{
			"pic/a/b/c",
			"pic/ttt.txt",
			"/pic",
			"pic",
			&mfs{},
			&dir{},
			nil,
		},
		{
			"pic/a/b/c",
			"pic/ttt.txt",
			"/a/v",
			"pic",
			&mfs{},
			&dir{},
			&os.PathError{},
		},
	}

	for i, v := range test {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if v.targetDir[0] == '/' {
				v.targetDir = v.targetDir[1:]
			}

			if !exists(v.dp) {
				if err := os.MkdirAll(v.dp, 0777); err != nil {
					t.Error(err)
				}
			}

			if !exists(v.df) {
				var f *os.File
				f, err := os.Create(v.df)
				if err != nil {
					t.Error(err)
				}
				f.Close()
			}

			_, err := Fork(v.targetDir)
			if err != nil && v.Err == nil {
				t.Error(err)
			}

			_ = os.RemoveAll(v.rmf)
		})
	}
}
