package main

import (
	"fmt"
	"io/ioutil"
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

func Thread(args ThrArgs) (num int, e error) {
	num = 0
	for {
		select {
		case fname := <-args.fnamechan:
			fmt.Printf("<%d>%s/%s %s/%s in %s\n", num, args.srcdir, fname, args.dstdir, fname, args.prefix)
			num++
		case ended := <-args.endchan:
			break
		}

	}

	args.closechan <- 1
	return

}

type PathFunc func(srcfile string, dstfile string, a ...interface{}) (num int, e error)

func PrintPath(srcfile string, dstfile string, a interface{}) (num int, e error) {
	var argptr []*ThrArgs
	if reflect.TypeOf(argptr) == reflect.TypeOf(a) {
		fmt.Printf("type assign\n")
	}

	return 0, nil
}

func PathDir(srcdir string, dstdir string, dep int, f PathFunc, a ...interface{}) int {
	var count int
	count = 0
	files, e := ioutil.ReadDir(srcdir)
	if e != nil {
		Debug("read directory %s error %v", srcdir, e)
		return count
	}

	for _, f := range files {
		if f.Mode().IsDir() {
			nsrcdir := srcdir + os.PathSeparator + f.Name()
			ndstdir := dstdir + os.PathSeparator + f.Name()
			nstat, e := os.Stat(ndstdir)
			if e != nil && os.IsNotExist(e) {
				os.Mkdir(ndstdir, nstat.Mode())
			}
			count += PathDir(nsrcdir, ndstdir, dep+1, f, a)
		} else {
			nsrcfile := srcdir + os.PathSeparator + f.Name()
			ndstfile := dstdir + os.PathSeparator + f.Name()
			f(nsrcfile, ndstfile, a)
			count++
		}
	}

	return count
}

func main() {
	var argsptr []*ThrArgs
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s srcdir dstdir\n", os.Args[0])
		os.Exit(2)
	}

	argsptr = make([]*ThrArgs, runtime.NumCPU())
	fnamechan := make(chan string)
	endchan := make(chan int)
	for i := 0; i < runtime.NumCPU(); i++ {
		prefix := fmt.Sprintf("Args%d", i)
		args := NewThrArgs(os.Args[1], os.Args[2], fnamechan, endchan, prefix)
		argsptr[i] = &args
		go Thread(args)
	}

}
