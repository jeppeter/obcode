package main

import (
	"fmt"
	"os"
	"runtime"
)

func LogFormat(format string, a ...interface{}) {

	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	fmt.Printf(s)
}

func main() {
	LogFormat("hello %s ss %s\n", os.Args[1], os.Args[0])
}
