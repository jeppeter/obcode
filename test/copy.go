package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

func Error(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stderr, s)
	return len(s)
}

func CopyFileContents(srcfile string, dstfile string) error {
	in, err := os.Open(srcfile)
	if err != nil {
		Error("could not open %s <%v>", srcfile, err)
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dstfile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Error("could not open %s for write <%v>", dstfile, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		Error("could not copy <%s => %s> <%v>", srcfile, dstfile, err)
		return err
	}

	err = out.Sync()
	if err != nil {
		Error("could not sync %s <%v>", dstfile, err)
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s srcfile dstfile\n", os.Args[0])
		os.Exit(3)
	}

	CopyFileContents(os.Args[1], os.Args[2])
}
