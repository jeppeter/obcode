package main

import (
	"fmt"
	"os"
)

func main() {
	for _, f := range os.Args[1:] {
		fileinfo, e := os.Lstat(f)
		if e != nil {
			fmt.Fprintf(os.Stderr, "can not stat %s\n", f)
		} else {
			fmt.Fprintf(os.Stdout, "%s perm 0%o\n", f, fileinfo.Mode().Perm())
			if fileinfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				l, _ := os.Readlink(f)
				fmt.Fprintf(os.Stdout, "%s is link %s\n", f, l)
			}
		}
	}
	return
}
