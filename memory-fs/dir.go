package mfs

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type dir struct {
	sync.RWMutex
	name    string
	perm    fs.FileMode
	modTime time.Time
	node    map[string]openFS
}

type openDir struct {
	*dir
	n     int
	entry []string
}

func (d *dir) FS() fs.File                { return &openDir{dir: d} }
func (d *dir) Name() string               { return d.name }
func (d *dir) Size() int64                { return 0 }
func (d *dir) Mode() fs.FileMode          { return d.perm }
func (d *dir) Type() fs.FileMode          { return d.perm }
func (d *dir) ModTime() time.Time         { return d.modTime }
func (d *dir) IsDir() bool                { return true }
func (d *dir) Sys() any                   { return nil }
func (d *dir) Stat() (fs.FileInfo, error) { return d, nil }
func (d *dir) Info() (fs.FileInfo, error) { return d, nil }
func (d *dir) Close() error               { return nil }
func (d *dir) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.name, Err: errors.New("is a directory")}
}
func (d *openDir) ReadDir(count int) ([]fs.DirEntry, error) {
	if d.n == 0 {
		d.RLock()
		if length := len(d.node); length != 0 {
			d.entry = make([]string, 0, len(d.node))

			for k := range d.node {
				d.entry = append(d.entry, k)
			}
		}
		d.RUnlock()

		sort.Strings(d.entry)
	}

	n := len(d.entry) - d.n
	if n == 0 {
		if count <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}
	if count > 0 && n > count {
		n = count
	}

	list := make([]fs.DirEntry, n)
	d.RLock()
	for i := range list {
		list[i] = d.node[d.entry[i+d.n]]
	}
	d.RUnlock()
	d.n += n
	return list, nil
}

var (
	_ openFS         = (*dir)(nil)
	_ fs.File        = (*openDir)(nil)
	_ fs.ReadDirFile = (*openDir)(nil)
)

func (d *dir) find(name string, idx int) (openFS, error) {
	i := strings.IndexByte(name[idx:], '/')
	switch i {
	case -1:
		return d.findNode(name[idx:])
	case 0:
		return d, nil
	}

	i += idx
	node, err := d.findNode(name[idx:i])
	if err != nil {
		return nil, err
	}
	if !node.IsDir() {
		return nil, fs.ErrNotExist
	}
	return node.(*dir).find(name, i+1)
}

func (d *dir) findNode(name string) (openFS, error) {
	d.RLock()
	node, ok := d.node[name]
	d.RUnlock()

	if !ok {
		return nil, fs.ErrNotExist
	}
	return node, nil
}

func (d *dir) writeFile(name string, data []byte, perm os.FileMode) (err error) {
	d.Lock()
	if _, ok := d.node[name]; ok {
		err = fs.ErrExist
	} else {
		d.node[name] = &file{
			name:    name,
			perm:    perm,
			modTime: time.Now(),
			data:    data,
		}
	}
	d.Unlock()

	return
}

func (d *dir) deleteNode(name string) (err error) {
	d.Lock()
	if _, ok := d.node[name]; !ok {
		err = fs.ErrNotExist
	} else {
		delete(d.node, name)
	}
	d.Unlock()
	return
}
