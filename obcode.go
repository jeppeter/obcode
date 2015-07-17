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

func Obcode(srcdir string, dstdir string, fname string) (replc int, e error) {
	sfile := srcdir + string(os.PathSeparator) + fname
	rfile, e := os.Open(sfile)
	if e != nil {
		Debug("can not open %s error %v", sfile, e)
		return 0, e
	}
	defer rfile.Close()
	dfile := dstdir + string(os.PathSeparator) + fname
	wfile, e := os.Create(dfile)
	if e != nil {
		Debug("can not open %s for writing error %v", dfile, e)
		return 0, e
	}
	defer wfile.Close()

	/************************************
	*          to read every one line ,and
	*          we write the line
	*
	************************************/
}

/********************************************
*    ch for the string to handle
*    done for the over down it
********************************************/
func ObcodeRoutine(srcdir string, dstdir string, ch chan string, done chan int, over chan int, idxnum int) {
	var fname string

	for {
		select {
		case fname := <-ch:
			Obcode(srcdir, dstdir, fname)
		case <-done:
			goto out_chan

		}
	}
out_chan:
	over <- 1
	return
}
