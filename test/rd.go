package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func PathDir(d string, dep int) int {
	var count int
	count = 0
	files, e := ioutil.ReadDir(d)

	if e != nil {
		fmt.Fprintf(os.Stderr, "%s read error %v\n", d, e)
		return 0
	}

	for i, f := range files {
		if f.Mode().IsDir() {
			fmt.Fprintf(os.Stdout, "%d[%d] %s%s%s\n", dep, i, d, string(os.PathSeparator), f.Name())
			count += PathDir(d+string(os.PathSeparator)+f.Name(), dep+1)
		} else {
			fmt.Fprintf(os.Stdout, "%d<%d> %s%s%s\n", dep, i, d, string(os.PathSeparator), f.Name())
			count++
		}
	}

	return count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s directory\n", os.Args[0])
		os.Exit(1)
	}
	cnt := PathDir(os.Args[1], 0)
	fmt.Fprintf(os.Stdout, "<%s> total count %d\n", os.Args[1], cnt)
	return
}
