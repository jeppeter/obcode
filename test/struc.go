package main

import (
	"fmt"
	"runtime"
)

type ThrArgs struct {
	srcdir string
	dstdir string
}

func main() {
	var argsarr []*ThrArgs

	argsarr = make([]*ThrArgs, runtime.NumCPU())
	fmt.Printf("hello\n")
	for i := 0; i < runtime.NumCPU(); i++ {
		argsarr[i] = &ThrArgs{}
		argsarr[i].srcdir = "hello"
		argsarr[i].dstdir = "hed"
	}
}
