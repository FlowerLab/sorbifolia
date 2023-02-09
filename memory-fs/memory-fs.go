package mfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
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
			name:    "/",
			modTime: time.Now(),
			node:    make(map[string]openFS),
		},
	}
}

func Fork(name string) (MemoryFS, error) {
	return fork("/", name)
}

func fork(memDirname, name string) (MemoryFS, error) {
	files, err := os.ReadDir(name)
	if err != nil {
		return nil, err
	}

	if name[len(name)-1] != '/' {
		name += "/"
	}

	root := &dir{
		name:    memDirname,
		modTime: time.Now(),
		node:    make(map[string]openFS),
	}

	for _, f := range files {
		tf := fmt.Sprintf("%s%s", name, f.Name())

		if !f.IsDir() {
			var content []byte
			content, err = os.ReadFile(tf)
			if err != nil {
				return nil, err
			}

			root.node[f.Name()] = &file{
				name:    f.Name(),
				modTime: time.Now(),
				data:    content,
			}
			continue
		}

		var cdir MemoryFS
		cdir, err = fork(f.Name(), tf)
		if err != nil {
			return nil, err
		}
		root.node[f.Name()] = cdir.(*mfs).root
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
