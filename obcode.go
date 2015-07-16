package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

func Obcode(srcdir string, dstdir string, fname string) (replc int, e error) {
	sfile := srcdir + string(os.PathSeparator) + fname
	rfile, e := os.Open(sfile)
	if e != nil {
		return e
	}
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
