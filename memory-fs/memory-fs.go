package mfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"
	"time"
)

type MemoryFS interface {
	fs.FS

	WriteFile(path string, data []byte, perm os.FileMode) error
	Mkdir(name string) error
	MkdirAll(name string) error
	Remove(name string) error
	Copy(name, to string) error
	Move(name, to string) error
}

type openFS interface {
	fs.DirEntry
	FS() fs.File
}

func New() MemoryFS {
	return &mfs{
		root: &dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			modTime: time.Now(),
			node:    make(map[string]openFS),
		},
	}
}

func Fork(name string) (MemoryFS, error) {
	_, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if name[0] == '/' {
		name = name[1:]
	}

	paths := strings.Split(name, "/")
	root := &dir{
		RWMutex: sync.RWMutex{},
		name:    "/",
		modTime: time.Now(),
		node:    make(map[string]openFS),
	}

	curr := root
	for i := 0; i < len(paths); i++ {
		nd := &dir{
			RWMutex: sync.RWMutex{},
			name:    paths[i],
			perm:    curr.perm,
			modTime: time.Now(),
			node:    make(map[string]openFS),
		}
		curr.node[paths[i]] = nd
		curr = nd
	}

	return &mfs{root}, nil
}

func Persistence(m MemoryFS, name string) error {
	f, err := os.Stat(name)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is not a directory", name)
	}

	d := m.(*mfs).root
	rec := make(map[string][]openFS)
	rec[name] = append(rec[name], openFS(d))
	// BFS
	for len(rec) > 0 {
		tmp := rec
		rec = make(map[string][]openFS)
		for k, v := range tmp {
			for _, ofs := range v {
				var cpath string
				if ofs.Name() != "/" {
					if k == "/" {
						cpath = fmt.Sprintf("/%s", ofs.Name())
					} else {
						cpath = fmt.Sprintf("%s/%s", k, ofs.Name())
					}
				}

				if !ofs.IsDir() {
					mf := ofs.(*file)

					var cf *os.File
					cf, err = os.Create(cpath)
					if err != nil {
						return err
					}
					_, _ = cf.Write(mf.data)
					_ = cf.Close()
					continue
				}
				md := ofs.(*dir)
				if md.name != "/" {
					_ = os.Mkdir(cpath, md.perm)
				}

				if len(cpath) == 0 {
					cpath = k
				}
				for _, mdn := range md.node {
					rec[cpath] = append(rec[cpath], mdn)
				}
			}
		}
	}

	return nil
}

var errIsDirectory = errors.New("is a directory")
