package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

func Debug(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stdout, s)
	return len(s)
}

func Error(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stderr, s)
	return len(s)
}

func Obcode(srcdir string, dstdir string, fname string, prefix string) (replc int, e error) {
	sfile := srcdir + string(os.PathSeparator) + fname
	dfile := dstdir + string(os.PathSeparator) + fname
	return ReadWriteFile(sfile, dfile, prefix)
}

/********************************************
*    ch for the string to handle
*    done for the over down it
********************************************/
func ObcodeRoutine(srcdir string, dstdir string, ch chan string, done chan int, over chan int, prefix string) (cnt int, e error) {
	var fname string
	cnt = 0
	e = nil
	for {
		select {
		case fname := <-ch:
			repl, err := Obcode(srcdir, dstdir, fname, prefix)
			if err != nil {
				Error("Obcode<%s%c%s>  in %s error %v", srcdir, os.PathSeparator, fname, prefix, err)
				e = err
				goto out_chan
			} else {
				cnt++
			}
		case <-done:
			goto out_chan

		}
	}
out_chan:
	over <- 1
	return cnt, e
}

func MainDispatch(srcdir string, dstdir string, partfile string, ch chan string) (e error) {
	var sd, dd, curs, curd string

	if len(partfile) > 0 {
		sd = srcdir + os.PathSeparator + partfile
		dd = dstdir + os.PathSeparator + partfile
	} else {
		sd = srcdir
		dd = dstdir
	}

	files, e := ioutil.ReadDir(sd)
	for i, f := range files {
		if f.Mode().IsDir() {
			curd = os + dd.PathSeparator + f.Name()
			os.Mkdir(name, perm)
		}
	}
}
