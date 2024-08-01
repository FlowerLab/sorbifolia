package debugger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

func GoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	field := strings.Fields(string(bytes.TrimPrefix(b, []byte("goroutine "))))
	n, _ := strconv.ParseUint(field[0], 10, 64)
	return n
}

func Ax() {
	fmt.Println(GoroutineID())
}

func CallStack(skip int) *runtime.Frames {
	skip += 2
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])

	return runtime.CallersFrames(pcs[:n])
}

func HeapDumpToFile(filename string) (_ string, err error) {
	var file *os.File
	if filename == "" {
		file, err = os.CreateTemp("", "heap-dump-*")
	} else {
		file, err = os.Create(filename)
	}
	if err != nil {
		return "", err
	}

	debug.WriteHeapDump(file.Fd())
	return file.Name(), file.Close()
}
