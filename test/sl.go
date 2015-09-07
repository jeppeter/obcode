package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s lnfile bfile\n", os.Args[0])
		os.Exit(3)
	}

	e := os.Symlink(os.Args[2], os.Args[1])
	if e != nil {
		fmt.Fprintf(os.Stderr, "can not link %s to %s error %v\n", os.Args[1], os.Args[2], e)
		os.Exit(4)
	}
	return
}
