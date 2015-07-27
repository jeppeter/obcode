package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"math/rand"
)

func Debug(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stdout, s)
	return len(s)
}

func ObFuncVarTransition(line string, funcname string, prefix string) retstr string {
	var rndstr string
	var rndbytes [20]byte
	for i:= 0 ;i < 20 ;i ++{
		rndnum := int(rand.Float32()* 100)
		rndnum %= 62
		if rndnum < 26 {
			rndbyte[i] = byte(int('a') + rndnum)
		} else if rndnum < 52 {
			rndbyte[i] = byte(int('A') + rndnum - 26)
		} else {
			rndbyte[i] = byte(int('0') + rndnum - 52)
		}
	}
	rndstr = string(rndbyte[:])
	retstr = fmt.Sprintf("#define %s %s_%s", funcname,prefix,rndstr)
	return retstr
}


func ObCodeTransition(line string,codename string,prefix string) retstr string {
	var 
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
	linenum := 1
	obfunc_reg, e := regexp.Compile(`OB_FUNC(\s+)([^(]+)`)
	if e != nil {
		Debug("could not make OB_FUNC regexp %v", e)
		return 0, e
	}
	obvar_reg, e := regexp.Compile(`OB_VAR(\s+)([^ \t;,=]+)`)
	if e != nil {
		Debug("could not make OB_VAR regexp %v", e)
		return 0, e
	}
	obcode_reg, e := regexp.Compile(`OB_CODE\(([^)]+)\)`)
	if e != nil {
		Debug("could not make OB_CODE regex %v", e)
		return 0, e
	}
	define_reg, e := regexp.Compile(`#define`)
	if e != nil {
		Debug("could not make define regexp %v", e)
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

		if obfunc_reg.Match(line) {
			r := obfunc_reg.FindStringSubmatch(string(line))
			Debug("<%d> func %s", linenum, r[2])
			if !define_reg.Match(string(line)) {

			}
		} else if obvar_reg.Match(line) {
			r := obvar_reg.FindStringSubmatch(string(line))
			Debug("<%d> var %s", linenum, r[2])
			if !define_reg.Match(string(line)) {

			}
		} else if obcode_reg.Match(line) {
			r := obcode_reg.FindStringSubmatch(string(line))
			Debug("<%d> code var %s", linenum, r[1])
			if !define_reg.Match(string(line)) {

			}
		}
		linenum++

	}
	return repl, nil
}

func main() {
	ReadWriteFile(os.Args[1], os.Args[2])
}
