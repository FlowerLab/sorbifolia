package mfs

import (
	"errors"
	"io/fs"
	"os"
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

func New() MemoryFS                               { panic("not implemented") }
func Fork(name string) (MemoryFS, error)          { panic("not implemented") }
func Persistence(mfs MemoryFS, name string) error { panic("not implemented") }

var errIsDirectory = errors.New("is a directory")
