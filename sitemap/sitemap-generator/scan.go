package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func findFile(path, dir string) []string {
	var (
		files, err = ioutil.ReadDir(filepath.Join(path, dir))
		arr        []string
	)
	if err != nil {
		return nil
	}

	for _, v := range files {
		if v.IsDir() {
			arr = append(arr, findFile(path, filepath.Join(dir, v.Name()))...)
			continue
		}
		if strings.HasSuffix(v.Name(), ".html") {
			arr = append(arr, filepath.Join(dir, v.Name()))
		}
	}

	return arr
}
