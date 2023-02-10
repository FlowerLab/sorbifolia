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

	return persistence(m, name)
}

func persistence(m MemoryFS, diskName string) error {
	d := m.(*mfs).root
	if diskName[len(diskName)-1] != '/' {
		diskName += "/"
	}

	for _, f := range d.node {
		cpath := fmt.Sprintf("%s%s", diskName, f.Name())
		if !f.IsDir() {
			mf := f.(*file)
			var cf *os.File
			cf, err := os.Create(cpath)
			if err != nil {
				return err
			}
			_, _ = cf.Write(mf.data)
			_ = cf.Close()
			continue
		}
		_ = os.Mkdir(cpath, os.ModePerm)

		if err := persistence(&mfs{f.(*dir)}, cpath); err != nil {
			return err
		}
	}
	return nil
}

var errIsDirectory = errors.New("is a directory")
