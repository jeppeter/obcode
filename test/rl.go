package main

import (
	"bufio"
	"fmt"
	"os"
)

func HandleFile(fname string) error {
	fh, e := os.Open(fname)
	if e != nil {
		fmt.Fprintf(os.Stderr, "can not open %s error %v\n", fname, e)
		return e
	}
	defer fh.Close()
	bufscan := bufio.NewScanner(fh)
	bufscan.Split(bufio.ScanLines)
	i := 0

	for bufscan.Scan() {
		t := bufscan.Text()
		fmt.Fprintf(os.Stdout, "<%d>%s\n", i, t)
		i++
	}
	return nil
}

func main() {
	var f string

	for _, f = range os.Args[1:] {
		HandleFile(f)
	}
}
