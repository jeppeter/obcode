package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s exp string", os.Args[0])
		os.Exit(1)
	}
	re, e := regexp.Compile(os.Args[1])
	if e != nil {
		fmt.Fprintf(os.Stderr, "compile %s error %v\n", os.Args[1], e)
		os.Exit(2)
	}

	buf := []byte(os.Args[2])
	b := re.Match(buf)

	if b {
		fmt.Fprintf(os.Stdout, "<%s> match <%s>\n", os.Args[1], os.Args[2])
	} else {
		fmt.Fprintf(os.Stdout, "(%s) not match (%s)\n", os.Args[1], os.Args[2])
	}
	return
}
