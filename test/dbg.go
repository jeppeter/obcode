package main

import (
	"fmt"
	"os"
	"runtime"
)

func LogFormat(format string, a ...interface{}) {

	_, f, l, _ := runtime.Caller(1)
	fmt.Printf("[%s:%d]\t", f, l)
	fmt.Printf(format, a...)
}

func main() {
	LogFormat("hello %s ss %s\n", os.Args[1], os.Args[0])
}
