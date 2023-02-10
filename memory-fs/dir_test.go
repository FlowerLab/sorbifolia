package mfs

import (
	"fmt"
	"io"
	"io/fs"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		openDir
		res []string
		rdc []int
		idx int
		Err error
	}{
		{
			openDir{
				dir: &dir{
					RWMutex: sync.RWMutex{},
					name:    "/",
					modTime: time.Now(),
					node: map[string]openFS{
						"a": &dir{name: "A"},
						"b": &dir{name: "B"},
						"c": &dir{name: "C"},
						"d": &dir{name: "D"},
						"f": &dir{name: "F"},
					},
				},
			},
			[]string{"A", "B", "C", "D", "F"},
			[]int{3, 2},
			0,
			nil,
		},
		{
			openDir{
				dir: &dir{
					RWMutex: sync.RWMutex{},
					name:    "/",
					modTime: time.Now(),
					node:    make(map[string]openFS),
				},
			},
			nil,
			[]int{-1},
			0,
			nil,
		},
		{
			openDir{
				dir: &dir{
					RWMutex: sync.RWMutex{},
					name:    "/",
					modTime: time.Now(),
					node:    make(map[string]openFS),
				},
			},
			nil,
			[]int{2},
			0,
			io.EOF,
		},
		{
			openDir{
				dir: &dir{
					RWMutex: sync.RWMutex{},
					name:    "/",
					modTime: time.Now(),
					node: map[string]openFS{
						"a": &dir{name: "A"},
						"b": &dir{name: "B"},
						"c": &dir{name: "C"},
						"d": &dir{name: "D"},
						"f": &dir{name: "F"},
					},
				},
			},
			[]string{"A", "B", "C", "D", "F"},
			[]int{3, -1},
			0,
			nil,
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			for _, c := range v.rdc {
				f, err := v.ReadDir(c)
				if !reflect.DeepEqual(v.Err, err) {
					t.Errorf("want %v, get %v", v.Err, err)
				}
				if v.res == nil && f != nil {
					t.Error("count < 0 is error")
				}
				for _, e := range f {
					if !reflect.DeepEqual(v.res[v.idx], e.Name()) {
						t.Errorf("want %v, get %v", v.res[v.idx], e.Name())
					}
					v.idx++
				}
			}
		})
	}
}

func TestFind(t *testing.T) {
	// /a/b/c/d 1
	tests := []struct {
		*mfs
		path   string
		idx    int
		expect string
		Err    error
		fp     string
	}{
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				modTime: time.Time{},
				node:    make(map[string]openFS),
			}},
			"/a/b/c",
			1,
			"c",
			nil,
			"/a/b/c",
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				modTime: time.Time{},
				node:    make(map[string]openFS),
			}},
			"b/c/d",
			0,
			"d",
			nil,
			"b/c/d",
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				modTime: time.Time{},
				node:    make(map[string]openFS),
			}},
			"b/c/d",
			0,
			"d",
			fs.ErrNotExist,
			"b/c/a",
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				modTime: time.Time{},
				node:    make(map[string]openFS),
			}},
			"b/c/d",
			0,
			"d",
			fs.ErrNotExist,
			"b/c/d/a/f",
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if err := v.MkdirAll(v.path[v.idx:]); err != nil {
				t.Error(err)
			}
			fd, err := v.root.find(v.fp, v.idx)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect %v,get %v", v.Err, err)
			}
			if fd != nil && !reflect.DeepEqual(v.expect, fd.Name()) {
				t.Errorf("expect %v,get %v", v.expect, fd.Name())
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		*dir
		res  []byte
		Err  error
		fn   string
		data []byte
	}{
		{
			&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				node: map[string]openFS{
					"a.txt": &file{
						name: "a.txt",
						data: nil,
					},
				},
			},
			nil,
			fs.ErrExist,
			"a.txt",
			nil,
		},
		{
			&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				node:    make(map[string]openFS),
			},
			[]byte{1, 2, 3, 4, 5},
			nil,
			"a.txt",
			[]byte{1, 2, 3, 4, 5},
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := v.writeFile(v.fn, v.data, v.perm)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
			if err == nil && !reflect.DeepEqual(v.res, v.node[v.fn].(*file).data) {
				t.Errorf("expect: %v,get: %v", v.res, v.node[v.fn].(*file).data)
			}
		})
	}
}

func TestDeleteNode(t *testing.T) {
	tests := []struct {
		*dir
		Err error
	}{
		{
			&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				node: map[string]openFS{
					"a.txt": &file{
						name: "a.txt",
						data: nil,
					},
				},
			},
			nil,
		},
		{
			&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				node:    make(map[string]openFS),
			},
			fs.ErrNotExist,
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := v.deleteNode("a.txt")
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
		})
	}
}

func TestAttribute(t *testing.T) {
	ct := time.Now()
	test := []struct {
		*dir
		size int64
		m    fs.FileMode
		tp   fs.FileMode
		ct   time.Time
		sys  any
	}{
		{
			&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				modTime: ct,
				node:    nil,
			},
			0,
			0,
			0,
			ct,
			nil,
		},
	}

	for i, v := range test {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if !reflect.DeepEqual(v.Size(), v.size) {
				t.Errorf("expect: %v,get: %v", v.Size(), v.size)
			}
			if !reflect.DeepEqual(v.Mode(), v.m) {
				t.Errorf("expect: %v,get: %v", v.Mode(), v.m)
			}
			if !reflect.DeepEqual(v.Type(), v.tp) {
				t.Errorf("expect: %v,get: %v", v.Type(), v.tp)
			}
			if !reflect.DeepEqual(v.ModTime(), v.modTime) {
				t.Errorf("expect: %v,get: %v", v.ModTime(), v.modTime)
			}
			if !reflect.DeepEqual(v.Sys(), v.sys) {
				t.Errorf("expect: %v,get: %v", v.Size(), v.size)
			}
		})
	}
}
