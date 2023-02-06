package mfs

import (
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestFsOpen(t *testing.T) {
	tests := []struct {
		*mfs
		Err    error
		expect string
		fn     string
	}{
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				modTime: time.Now(),
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						modTime: time.Now(),
						node: map[string]openFS{
							"b": &dir{
								RWMutex: sync.RWMutex{},
								name:    "b",
								modTime: time.Now(),
								node: map[string]openFS{
									"c.txt": &file{
										name:    "c.txt",
										perm:    0,
										modTime: time.Now(),
										data:    []byte{1, 2, 3, 4, 5},
									},
								},
							},
						},
					},
				},
			}},
			nil,
			"c.txt",
			"a/b/c.txt",
		},
		{
			&mfs{},
			fs.ErrInvalid,
			"",
			"/a/b/c.txt",
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				modTime: time.Now(),
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						modTime: time.Now(),
						node: map[string]openFS{
							"b": &dir{
								RWMutex: sync.RWMutex{},
								name:    "b",
								modTime: time.Now(),
								node: map[string]openFS{
									"c.txt": &file{
										name:    "c.txt",
										perm:    0,
										modTime: time.Now(),
										data:    []byte{1, 2, 3, 4, 5},
									},
								},
							},
						},
					},
				},
			}},
			fs.ErrNotExist,
			"",
			"a/b.txt",
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			f, err := v.Open(v.fn)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
			if err == nil {
				var tf fs.FileInfo
				tf, err = f.Stat()
				if err != nil {
					t.Log(err)
				}
				if !reflect.DeepEqual(v.expect, tf.Name()) {
					t.Errorf("expect: %v,get: %v", v.expect, tf.Name())
				}
			}
		})
	}
}

func TestFsWriteFile(t *testing.T) {
	tests := []struct {
		*mfs
		path string
		data []byte
		perm os.FileMode
		Err  error
		idx  int
	}{
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				modTime: time.Now(),
				node:    make(map[string]openFS),
			}},
			"b.txt",
			[]byte("123456"),
			0,
			nil,
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				modTime: time.Now(),
				node: map[string]openFS{
					"b.txt": &file{
						name:    "b.txt",
						perm:    0,
						modTime: time.Now(),
						data:    []byte("123"),
					},
				},
			}},
			"b.txt",
			[]byte("123456"),
			0,
			fs.ErrExist,
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				modTime: time.Now(),
				node:    make(map[string]openFS),
			}},
			"/b.txt",
			[]byte("123456"),
			0,
			nil,
			1,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				modTime: time.Time{},
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						perm:    0,
						modTime: time.Now(),
						node: map[string]openFS{
							"c": &dir{
								RWMutex: sync.RWMutex{},
								name:    "c",
								perm:    0,
								modTime: time.Now(),
								node:    make(map[string]openFS),
							},
						},
					},
				},
			}},
			"a/c/b.txt",
			[]byte("123456"),
			0,
			nil,
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				modTime: time.Time{},
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						perm:    0,
						modTime: time.Now(),
						node: map[string]openFS{
							"c": &dir{
								RWMutex: sync.RWMutex{},
								name:    "c",
								perm:    0,
								modTime: time.Now(),
								node:    make(map[string]openFS),
							},
						},
					},
				},
			}},
			"/a/c/b.txt",
			[]byte("123456"),
			0,
			nil,
			1,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				modTime: time.Time{},
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						perm:    0,
						modTime: time.Now(),
						node: map[string]openFS{
							"c": &dir{
								RWMutex: sync.RWMutex{},
								name:    "c",
								perm:    0,
								modTime: time.Now(),
								node:    make(map[string]openFS),
							},
						},
					},
				},
			}},
			"/a/c/d/b.txt",
			[]byte("123456"),
			0,
			fs.ErrNotExist,
			1,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				perm:    0,
				modTime: time.Time{},
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						perm:    0,
						modTime: time.Now(),
						node: map[string]openFS{
							"c": &file{
								name: "c",
								perm: 0,
								data: []byte("dd"),
							},
						},
					},
				},
			}},
			"/a/c/b.txt",
			[]byte("123456"),
			0,
			&fs.PathError{
				Op:   "write",
				Path: "/a/c/b.txt",
				Err:  fmt.Errorf("%s isn't a directory", "/a/c"),
			},
			1,
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := v.WriteFile(v.path, v.data, v.perm)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
			if err == nil {
				var (
					tf   fs.FileInfo
					find openFS
				)
				find, err = v.root.find(v.path, v.idx)
				if err != nil {
					t.Error(err)
				}
				tf, err = find.FS().Stat()
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(v.data, tf.(*file).data) {
					t.Errorf("expect: %v,get: %v", v.data, tf.(*file).data)
				}
			}
		})
	}
}

func TestFsRemove(t *testing.T) {
	tests := []struct {
		*mfs
		path string
		Err  error
	}{
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"b.txt": &file{
						name: "b.txt",
						data: nil,
					},
				},
			}},
			"a/b.txt",
			fs.ErrNotExist,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"b.txt": &file{
						name: "b.txt",
						data: nil,
					},
				},
			}},
			"b.txt",
			nil,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"b.txt": &file{
						name: "b.txt",
						data: nil,
					},
				},
			}},
			"/b.txt",
			nil,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						node: map[string]openFS{
							"d": &dir{
								RWMutex: sync.RWMutex{},
								name:    "d",
								node: map[string]openFS{
									"b": &dir{
										RWMutex: sync.RWMutex{},
										name:    "b",
										node:    nil,
									},
								},
							},
						},
					},
				},
			}},
			"/a/d/b",
			nil,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						node: map[string]openFS{
							"d": &file{
								name: "d",
								data: nil,
							},
						},
					},
				},
			}},
			"/a/d/b",
			&fs.PathError{
				Op:   "delete",
				Path: "/a/d/b",
				Err:  fmt.Errorf("%s isn't a directory", "/a/d"),
			},
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := v.Remove(v.path)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
		})
	}
}

func TestFsCopy(t *testing.T) {
	tests := []struct {
		*mfs
		name string
		to   string
		Err  error
		res  []string
		idx  int
	}{
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &file{
						name: "a",
						perm: 0,
						data: []byte("dsa"),
					},
					"b": &dir{
						RWMutex: sync.RWMutex{},
						name:    "b",
						node:    make(map[string]openFS),
					},
				},
			}},
			"a",
			"b/c",
			nil,
			[]string{"b/c"},
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						node: map[string]openFS{
							"z": &dir{
								RWMutex: sync.RWMutex{},
								name:    "z",
								node:    make(map[string]openFS),
							},
							"x": &file{
								name: "x",
								data: nil,
							},
							"y": &dir{
								RWMutex: sync.RWMutex{},
								name:    "y",
								node:    make(map[string]openFS),
							},
						},
					},
					"b": &dir{
						RWMutex: sync.RWMutex{},
						name:    "b",
						node:    make(map[string]openFS),
					},
				},
			}},
			"a",
			"b",
			fs.ErrExist,
			nil,
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						node: map[string]openFS{
							"z": &dir{
								RWMutex: sync.RWMutex{},
								name:    "z",
								node:    make(map[string]openFS),
							},
							"x": &file{
								name: "x",
								data: nil,
							},
							"y": &dir{
								RWMutex: sync.RWMutex{},
								name:    "y",
								node:    make(map[string]openFS),
							},
						},
					},
					"b": &dir{
						RWMutex: sync.RWMutex{},
						name:    "b",
						node:    make(map[string]openFS),
					},
				},
			}},
			"a",
			"c",
			nil,
			[]string{"c/z", "c/x", "c/y"},
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						node: map[string]openFS{
							"z": &dir{
								RWMutex: sync.RWMutex{},
								name:    "z",
								node:    make(map[string]openFS),
							},
							"x": &file{
								name: "x",
								data: nil,
							},
							"y": &dir{
								RWMutex: sync.RWMutex{},
								name:    "y",
								node:    make(map[string]openFS),
							},
						},
					},
					"b": &dir{
						RWMutex: sync.RWMutex{},
						name:    "b",
						node: map[string]openFS{
							"c": &dir{
								RWMutex: sync.RWMutex{},
								name:    "c",
								perm:    0,
								node:    make(map[string]openFS),
							},
						},
					},
				},
			}},
			"a",
			"/b/c/d",
			nil,
			[]string{"b/c/d/z", "b/c/d/x", "b/c/d/y"},
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &file{
						name: "a",
						perm: 0,
						data: []byte("dsa"),
					},
					"b": &dir{
						RWMutex: sync.RWMutex{},
						name:    "b",
						node:    make(map[string]openFS),
					},
				},
			}},
			"a/d",
			"b/c",
			fs.ErrNotExist,
			[]string{"b/c"},
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						node:    make(map[string]openFS),
					},
					"b": &dir{
						RWMutex: sync.RWMutex{},
						name:    "b",
						node:    make(map[string]openFS),
					},
				},
			}},
			"a",
			"/c",
			nil,
			[]string{"c"},
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &file{
						name: "a",
						perm: 0,
						data: []byte("dsa"),
					},
					"b": &dir{
						RWMutex: sync.RWMutex{},
						name:    "b",
						node:    make(map[string]openFS),
					},
				},
			}},
			"b",
			"a/c",
			&fs.PathError{
				Op:   "delete",
				Path: "a/c",
				Err:  fmt.Errorf("%s isn't a directory", "a"),
			},
			[]string{"a/c"},
			0,
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						node:    make(map[string]openFS),
					},
					"b": &dir{
						RWMutex: sync.RWMutex{},
						name:    "b",
						node:    make(map[string]openFS),
					},
				},
			}},
			"a",
			"d/c",
			fs.ErrNotExist,
			[]string{"d/c"},
			0,
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := v.Copy(v.name, v.to)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
			if err == nil {
				for _, p := range v.res {
					if _, err = v.root.find(p, v.idx); err != nil {
						t.Error(err)
					}
				}
			}
		})
	}
}

func TestFsMove(t *testing.T) {
	tests := []struct {
		*mfs
		src string
		to  string
		Err error
		gs  string
	}{
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node: map[string]openFS{
				"a": &dir{
					RWMutex: sync.RWMutex{},
					name:    "a",
					node: map[string]openFS{
						"b": &file{
							name: "b",
							data: nil,
						},
					},
				},
				"c": &dir{
					RWMutex: sync.RWMutex{},
					name:    "c",
					node:    make(map[string]openFS),
				},
			},
		}},
			"a/b",
			"c/d",
			nil,
			"c/d/b",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node: map[string]openFS{
				"a": &dir{
					RWMutex: sync.RWMutex{},
					name:    "a",
					node: map[string]openFS{
						"b": &dir{
							RWMutex: sync.RWMutex{},
							name:    "b",
							node: map[string]openFS{
								"x": &dir{
									RWMutex: sync.RWMutex{},
									name:    "x",
									node:    nil,
								},
								"y": &file{
									name: "y",
									data: nil,
								},
								"z": &dir{
									RWMutex: sync.RWMutex{},
									name:    "z",
									node:    nil,
								},
							},
						},
					},
				},
				"c": &dir{
					RWMutex: sync.RWMutex{},
					name:    "c",
					node:    make(map[string]openFS),
				},
			},
		}},
			"a/b",
			"c/d",
			nil,
			"c/d/x",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node: map[string]openFS{
				"a": &dir{
					RWMutex: sync.RWMutex{},
					name:    "a",
					node: map[string]openFS{
						"b": &file{
							name: "b",
							data: nil,
						},
					},
				},
				"c": &dir{
					RWMutex: sync.RWMutex{},
					name:    "c",
					node:    make(map[string]openFS),
				},
			},
		}},
			"a/x",
			"c/d",
			fs.ErrNotExist,
			"",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node: map[string]openFS{
				"a": &dir{
					RWMutex: sync.RWMutex{},
					name:    "a",
					node: map[string]openFS{
						"b": &file{
							name: "b",
							data: nil,
						},
					},
				},
				"c": &dir{
					RWMutex: sync.RWMutex{},
					name:    "c",
					node:    make(map[string]openFS),
				},
			},
		}},
			"a/b",
			"d",
			nil,
			"d/b",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node: map[string]openFS{
				"a": &dir{
					RWMutex: sync.RWMutex{},
					name:    "a",
					node: map[string]openFS{
						"b": &file{
							name: "b",
							data: nil,
						},
					},
				},
				"c": &dir{
					RWMutex: sync.RWMutex{},
					name:    "c",
					node:    make(map[string]openFS),
				},
			},
		}},
			"a/b",
			"/d",
			nil,
			"d/b",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node: map[string]openFS{
				"a": &dir{
					RWMutex: sync.RWMutex{},
					name:    "a",
					node: map[string]openFS{
						"b": &file{
							name: "b",
							data: nil,
						},
					},
				},
				"c": &dir{
					RWMutex: sync.RWMutex{},
					name:    "c",
					node:    make(map[string]openFS),
				},
			},
		}},
			"a/b",
			"/c/x/y",
			fs.ErrNotExist,
			"",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node: map[string]openFS{
				"a": &dir{
					RWMutex: sync.RWMutex{},
					name:    "a",
					node: map[string]openFS{
						"b": &file{
							name: "b",
							data: nil,
						},
					},
				},
				"c": &dir{
					RWMutex: sync.RWMutex{},
					name:    "c",
					node: map[string]openFS{
						"x": &file{
							name: "x",
							data: nil,
						},
					},
				},
			},
		}},
			"a/b",
			"/c/x/y",
			&fs.PathError{
				Op:   "delete",
				Path: "/c/x/y",
				Err:  fmt.Errorf("%s isn't a directory", "/c/x"),
			},
			"",
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := v.Move(v.src, v.to)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
			if err == nil {
				if _, err = v.root.find(v.gs, 0); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestFsMkdirAll(t *testing.T) {
	tests := []struct {
		*mfs
		Err error
		dn  string
	}{
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node:    make(map[string]openFS),
			}},
			nil,
			"a/b/c",
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node:    make(map[string]openFS),
			}},
			nil,
			".",
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &dir{
						RWMutex: sync.RWMutex{},
						name:    "a",
						node: map[string]openFS{
							"b": &dir{
								RWMutex: sync.RWMutex{},
								name:    "b",
								node:    make(map[string]openFS),
							},
						},
					},
				},
			}},
			nil,
			"a/b/c/d",
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &file{
						name: "a",
						data: nil,
					},
				},
			}},
			fmt.Errorf("%s is not directory", "a"),
			"a/b/c/d",
		},
		{
			&mfs{&dir{
				RWMutex: sync.RWMutex{},
				name:    "/",
				node: map[string]openFS{
					"a": &file{
						name: "a",
						data: nil,
					},
				},
			}},
			fs.ErrInvalid,
			"/a/b/c/d",
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := v.MkdirAll(v.dn)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
		})
	}
}

func TestFsMkdir(t *testing.T) {
	tests := []struct {
		*mfs
		Err  error
		name string
	}{
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node:    make(map[string]openFS),
		}},
			nil,
			"abc",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node:    make(map[string]openFS),
		}},
			nil,
			"./abc",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node:    make(map[string]openFS),
		}},
			fs.ErrInvalid,
			"./a/bc",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node:    make(map[string]openFS),
		}},
			fs.ErrInvalid,
			"",
		},
		{&mfs{&dir{
			RWMutex: sync.RWMutex{},
			name:    "/",
			perm:    0,
			modTime: time.Time{},
			node: map[string]openFS{
				"abc": &dir{
					RWMutex: sync.RWMutex{},
					name:    "abc",
				},
			},
		}},
			fs.ErrExist,
			"abc",
		},
	}

	for i, v := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := v.Mkdir(v.name)
			if !reflect.DeepEqual(v.Err, err) {
				t.Errorf("expect: %v,get: %v", v.Err, err)
			}
		})
	}
}
