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

type mfs struct {
	root *dir // is root directory: /
}

func (m mfs) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
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

func (m mfs) Move(name, to string) error {
	f, err := m.Open(to)
	if err != nil {
		return err
	}

	if tf, ok := f.(*openFile); ok {
		var of fs.File
		of, err = m.Open(name)
		if err != nil {
			return err
		}
		if t, ok := of.(*openFile); ok {
			t.name = tf.name
		}
		if t, ok := of.(*openDir); ok {
			t.name = tf.name
		}
	} else {
		var (
			td  = f.(*openDir)
			idx = strings.LastIndexByte(name, '/')
			d   *dir
			tmp = name
		)

		switch idx {
		case -1:
			d = m.root
		case 0:
			d = m.root
			name = name[1:]
		default:
			name = name[idx+1:]
			var i int
			if tmp[0] == '/' {
				i = 1
			}

			var node openFS
			node, err = m.root.find(tmp[i:idx], i)
			if err != nil {
				return err
			}

			if !node.IsDir() {
				return fmt.Errorf("%s is not a directory", tmp[:idx])
			}
			d = node.(*dir)
		}

		d.Lock()
		td.Lock()

		md := d.node[name]
		td.node[name] = md
		err = d.deleteNode(name)

		td.Unlock()
		d.Lock()

		return err
	}

	return nil
}

func (m mfs) MkdirAll(name string) error {
	if !fs.ValidPath(name) {
		return errors.New("invalid path")
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
				RWMutex: sync.RWMutex{},
				name:    path,
				modTime: time.Now(),
				node:    make(map[string]openFS),
			}

			currentDir.Lock()
			currentDir.node[path] = childDir
			currentDir.Unlock()

			currentDir = childDir
		} else {
			if !tf.IsDir() {
				return fmt.Errorf("%s is not directory", path)
			}
			currentDir = tf.(*openDir).dir
		}
	}

	return nil
}

func (m mfs) Mkdir(name string) error {
	if len(name) == 0 {
		return fs.ErrInvalid
	}

	var dirname string
	paths := strings.Split(name, "/")
	if len(paths) == 1 {
		dirname = paths[0]
	} else if len(paths) == 2 && paths[0] == "." {
		dirname = paths[1]
	} else {
		return errors.New("the format of name is error")
	}

	m.root.Lock()
	defer m.root.Unlock()

	if _, ok := m.root.node[dirname]; !ok {
		childDir := &dir{
			RWMutex: sync.RWMutex{},
			name:    paths[0],
			modTime: time.Now(),
			node:    make(map[string]openFS),
		}
		m.root.node[paths[0]] = childDir
	} else {
		return fs.ErrExist
	}

	return nil
}
