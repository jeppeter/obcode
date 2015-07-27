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

func CountTabs(line string) (ntabs int) {
	var bs []byte
	var cnttabs, cntspcs int

	bs = []byte(line)
	cnttabs = 0
	cntspcs = 0
	for i := 0; i < len(bs); i++ {
		if bs[i] == '\t' {
			cnttabs++
		} else if bs[i] == ' ' {
			cntspcs++
		} else {
			break
		}
	}

	ntabs = 0
	ntabs += cnttabs
	ntabs += (cntspcs / 4)
	if (cntspcs % 4) != 0 {
		ntabs++
	}

	return ntabs
}
func ReadFileTabs(fname string) error {
	fh, e := os.Open(fname)
	if e != nil {
		Debug("Can not open %s error %v", fname, e)
		return e
	}

	defer fh.Close()

	rbuf := bufio.NewReaderSize(fh, 4096*4)

	linenum := 0
	for {
		line, _, e := rbuf.ReadLine()
		if e != nil {
			if e == io.EOF {
				break
			}
			Debug("can not read %s file", fname)
			return e
		}

		ntabs := CountTabs(string(line))
		Debug("<%d> tabs %d %s", linenum, ntabs, line)
		linenum++
	}

	return nil
}

func main() {
	ReadFileTabs(os.Args[1])
}
