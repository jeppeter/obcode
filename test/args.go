package main

import (
	"container/list"
	"fmt"
	"os"
	"reflect"
	"runtime"
)

type ThrArgs struct {
	fnamechan chan string
	endchan   chan int
	closechan chan int
	srcdir    string
	dstdir    string
	prefix    string
}

func Debug(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stdout, s)
	return len(s)
}

func NewThrArgs(srcdir string, dstdir string, fnamechan chan string, endchan chan int, prefix string) (args *ThrArgs) {
	args = new(ThrArgs)
	args.fnamechan = fnamechan
	args.endchan = endchan
	args.closechan = make(chan int)
	args.srcdir = srcdir
	args.dstdir = dstdir
	args.prefix = prefix
	return args
}

func Stringer(arg *ThrArgs) string {
	return fmt.Sprintf("prefix %s srcdir %s dstdir %s", arg.prefix, arg.srcdir, arg.dstdir)
}

type PathFunc func(srcfile string, dstfile string, a *list.List) (num int, e error)

func PrintPath(srcfile string, dstfile string, a *list.List) (num int, e error) {
	var argptr *ThrArgs
	for e := a.Front(); e != nil; e = e.Next() {
		if reflect.TypeOf(e.Value) == reflect.TypeOf(argptr) {
			argptr = e.Value.(*ThrArgs)
			fmt.Printf("%s\n", Stringer(argptr))
		}
	}

	return 0, nil
}

func main() {
	//var argsptr []*ThrArgs

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s srcdir dstdir\n", os.Args[0])
		os.Exit(2)
	}
	l := list.New()
	//argsptr = make([]*ThrArgs, runtime.NumCPU())
	fnamechan := make(chan string)
	endchan := make(chan int)
	for i := 0; i < runtime.NumCPU(); i++ {
		prefix := fmt.Sprintf("Args%d", i)
		args := NewThrArgs(os.Args[1], os.Args[2], fnamechan, endchan, prefix)
		l.PushBack(args)

	}

	PrintPath(os.Args[1], os.Args[2], l)

}
