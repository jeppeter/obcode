package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s string suffix", os.Args[0])
		os.Exit(3)
	}
	if strings.HasSuffix(os.Args[1], os.Args[2]) {
		fmt.Fprintf(os.Stdout, "%s has suffix %s\n", os.Args[1], os.Args[2])
	} else {
		fmt.Fprintf(os.Stdout, "%s not has suffix %s\n", os.Args[1], os.Args[2])
	}
	return
}
