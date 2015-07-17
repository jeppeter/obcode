package main

import (
	"bufio"
	"fmt"
	"io"
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

func ReadWriteFile(fname string, wfname string) (repl int, e error) {
	repl = 0
	rf, e := os.Open(fname)
	if e != nil {
		Debug("open %s error %v", fname, e)
		return 0, e
	}

	defer rf.Close()
	wf, e := os.OpenFile(wfname, os.O_WRONLY|os.O_CREATE, 0666)
	if e != nil {
		Debug("open write %s error %v", wfname, e)
		return 0, e
	}
	defer wf.Close()

	rbuf := bufio.NewReaderSize(rf, 4096*4)
	if e != nil {
		Debug("make reader size error %v", e)
		return 0, e
	}
	for {
		line, _, e := rbuf.ReadLine()
		if e != nil {
			if e == io.EOF {
				break
			}
			Debug("read %s line error %v", fname, e)
			return 0, e
		}
		wf.WriteString(string(line) + "\n")
	}
	return repl, nil
}

func main() {
	ReadWriteFile(os.Args[1], os.Args[2])
}
