package mfs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

type mfs struct {
	root *dir // is root directory: /
}

func (m mfs) Open(name string) (fs.File, error) {
	if len(name) == 0 {
		return nil, fs.ErrInvalid
	}

	var idx int
	if name[0] == '/' {
		idx = 1
	}
	node, err := m.root.find(name, idx)
	if err != nil {
		return nil, err
	}
	return node.FS(), nil
}

func (m mfs) WriteFile(path string, data []byte, perm os.FileMode) error {
	var (
		name = path
		i    = strings.LastIndexByte(name, '/')
		d    *dir
	)

	switch i {
	case -1:
		d = m.root
	case 0:
		name = name[1:]
		d = m.root
	default:
		name = name[i+1:]

		var idx int
		if path[0] == '/' {
			idx = 1
		}
		node, err := m.root.find(path[idx:i], idx)
		if err != nil {
			return err
		}
		if !node.IsDir() {
			return &fs.PathError{
				Op:   "write",
				Path: path,
				Err:  fmt.Errorf("%s isn't a directory", path[:i]),
			}
		}
		d = node.(*dir)
	}

	return d.writeFile(name, data, perm)
}

func (m mfs) Remove(path string) error {
	var (
		name = path
		i    = strings.LastIndexByte(name, '/')
		d    *dir
	)

	switch i {
	case -1:
		d = m.root
	case 0:
		name = name[1:]
		d = m.root
	default:
		name = name[i+1:]

		var idx int
		if path[0] == '/' {
			idx = 1
		}
		node, err := m.root.find(path[idx:i], idx)
		if err != nil {
			return err
		}
		if !node.IsDir() {
			return &fs.PathError{
				Op:   "delete",
				Path: path,
				Err:  fmt.Errorf("%s isn't a directory", path[:i]),
			}
		}
		d = node.(*dir)
	}

	return d.deleteNode(name)
}

func (m mfs) Copy(name, to string) error {
	f, err := m.Open(name)
	if err != nil {
		return err
	}

	if of, ok := f.(*openFile); ok {
		return m.WriteFile(to, of.data, of.perm)
	}

	od := f.(*openDir).dir
	name = to
	var (
		i = strings.LastIndexByte(name, '/')
		d *dir
	)

	switch i {
	case -1:
		d = m.root
	case 0:
		name = name[1:]
		d = m.root
	default:
		name = name[i+1:]

		var idx int
		if to[0] == '/' {
			idx = 1
		}
		var node openFS
		node, err = m.root.find(to[idx:i], idx)
		if err != nil {
			return err
		}
		if !node.IsDir() {
			return &fs.PathError{
				Op:   "delete",
				Path: to,
				Err:  fmt.Errorf("%s isn't a directory", to[:i]),
			}
		}
		d = node.(*dir)
	}

	nd := &dir{
		name:    od.name,
		perm:    od.perm,
		modTime: time.Now(),
		node:    make(map[string]openFS),
	}

	od.RLock()
	for k, v := range od.node {
		nd.node[k] = v
	}
	od.RUnlock()

	d.Lock()
	if _, ok := d.node[name]; ok {
		err = fs.ErrExist
	} else {
		d.node[name] = nd
	}
	d.Unlock()

	return err
}

func (m mfs) Move(name, to string) error { panic("not implemented") }
func (m mfs) MkdirAll(name string) error { panic("not implemented") }
func (m mfs) Mkdir(name string) error    { panic("not implemented") }
