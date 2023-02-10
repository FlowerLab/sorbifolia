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
	if !fs.ValidPath(name) {
		return nil, fs.ErrInvalid
	}

	node, err := m.root.find(name, 0)
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
		node, err := m.root.find(path[idx:i], 0)
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
		node, err := m.root.find(path[idx:i], 0)
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
		node, err = m.root.find(to[idx:i], 0)
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
		name:    name,
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

func (m mfs) Move(name, to string) error {
	f, err := m.Open(name)
	if err != nil {
		return err
	}

	source := name
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
		if node, err = m.root.find(to[idx:i], 0); err != nil {
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
		name:    name,
		perm:    d.perm,
		modTime: time.Now(),
		node:    make(map[string]openFS),
	}

	d.Lock()
	if _, ok := d.node[name]; ok {
		d.Unlock()
		return fs.ErrExist
	}
	d.node[name] = nd
	d.Unlock()

	if of, ok := f.(*openFile); ok {
		nf := &file{
			name:    of.name,
			perm:    of.perm,
			modTime: time.Now(),
			data:    of.data,
		}
		nd.node[nf.name] = nf
	} else {
		od := f.(*openDir).dir
		od.RLock()
		for k, v := range od.node {
			nd.node[k] = v
		}
		od.RUnlock()
	}

	return m.Remove(source)
}

func (m mfs) MkdirAll(name string) error {
	if !fs.ValidPath(name) {
		return fs.ErrInvalid
	}
	if name == "." {
		return nil
	}

	currentDir := m.root
	paths := strings.Split(name, "/")
	for _, path := range paths {
		currentDir.RLock()
		tf, ok := currentDir.node[path]
		currentDir.RUnlock()

		if !ok {
			childDir := &dir{
				name:    path,
				modTime: time.Now(),
				node:    make(map[string]openFS),
			}

			currentDir.Lock()
			currentDir.node[path] = childDir
			currentDir.Unlock()

			currentDir = childDir
			continue
		}

		if !tf.IsDir() {
			return fmt.Errorf("%s is not directory", path)
		}
		currentDir = tf.(*dir)
	}

	return nil
}

func (m mfs) Mkdir(name string) error {
	if len(name) == 0 || name[len(name)-1] == '/' {
		return fs.ErrInvalid
	}

	var (
		dirname string
		d       = m.root
		idx     = strings.LastIndexByte(name, '/')
	)

	switch idx {
	case 0:
		dirname = name[1:]
	case -1:
		dirname = name
	default:
		dirname = name[idx+1:]

		var i int
		if name[0] == '/' {
			i = 1
		}
		node, err := m.root.find(name[i:idx], 0)
		if err != nil {
			return err
		}
		if !node.IsDir() {
			return &fs.PathError{
				Op:   "delete",
				Path: name[i:idx],
				Err:  fmt.Errorf("%s isn't a directory", name[i:idx]),
			}
		}
		d = node.(*dir)
	}

	d.Lock()
	defer d.Unlock()

	if _, ok := d.node[dirname]; ok {
		return fs.ErrExist
	}
	d.node[dirname] = &dir{
		name:    dirname,
		perm:    d.perm,
		modTime: time.Now(),
		node:    make(map[string]openFS),
	}
	return nil
}
