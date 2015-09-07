package main

import (
	"fmt"
	"os"
	//"strings"
)

func main() {
	for _, s := range os.Args[1:] {
		if s == "hello" {
			fmt.Fprintf(os.Stdout, "%s hello\n", s)
		} else {
			fmt.Fprintf(os.Stdout, "%s not hello\n", s)
		}
	}
}
