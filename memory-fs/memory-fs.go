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
			if content, err = os.ReadFile(tf); err != nil {
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
		if cdir, err = fork(f.Name(), tf); err != nil {
			return nil, err
		}
		root.node[f.Name()] = cdir.(*mfs).root
	}

	return &mfs{root}, nil
}

func Persistence(m MemoryFS, name string) error {
	f, err := os.Stat(name)
	if err != nil {
		if !os.IsExist(err) {
			if err := os.MkdirAll(name, 0750); err != nil {
				return err
			}
		}
	}
	if err == nil && !f.IsDir() {
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
			if _, err = cf.Write(mf.data); err != nil {
				return err
			}
			if err = cf.Close(); err != nil {
				return err
			}
			continue
		}

		if !exists(cpath) {
			if err := os.Mkdir(cpath, os.ModePerm); err != nil {
				return err
			}
		}
		if err := persistence(&mfs{f.(*dir)}, cpath); err != nil {
			return err
		}
	}
	return nil
}

var errIsDirectory = errors.New("is a directory")

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
