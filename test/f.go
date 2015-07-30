package main

import (
	"fmt"
	"os"
)

type PathFunc func(srcfile string, dstfile string, a ...interface{}) (num int, e error)

func CopyFile(srcfile string, dstfile string, a ...interface{}) (num int, e error) {
	fmt.Printf("in copy file\n")
	return 0, nil
}

func PathDir(srcdir string, dstdir string, fn PathFunc, a ...interface{}) {
	fmt.Printf("srcdir %s dstdir %s\n", srcdir, dstdir)
	fn(srcdir, dstdir, a...)
}

func main() {
	PathDir(os.Args[1], os.Args[2], CopyFile, os.Args[3:])
}
