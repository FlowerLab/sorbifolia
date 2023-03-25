package main

import (
	"os"
	"os/exec"
)

func Compress(src, dst string) error {
	cmd := exec.Command("7z", "a", dst, "-m0=zstd", src)
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

func UnCompress(src, dst string) error {
	cmd := exec.Command("7z", "x", src, "-o="+dst)
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	return os.RemoveAll(src)
}
