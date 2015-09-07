package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

func ObFuncTransition(line string, funcname string, prefix string) (retstr string) {
	var rndstr string
	var rndbyte [20]byte
	for i := 0; i < 20; i++ {
		rndnum := int(rand.Float32() * 100)
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
	retstr = fmt.Sprintf("#define %s %s_%s\n", funcname, prefix, rndstr)
	return retstr
}

func ObVarTransition(line string, varname string, prefix string) (retstr string) {
	var rndstr string
	var rndbyte [20]byte
	for i := 0; i < 20; i++ {
		rndnum := int(rand.Float32() * 100)
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
	retstr = "#undef OB_VAR\n"
	retstr += fmt.Sprintf("#define OB_VAR(%s) %s_%s\n", varname, prefix, rndstr)
	retstr += fmt.Sprintf("#define %s %s_%s\n", varname, prefix, rndstr)
	retstr += line
	retstr += "\n"
	retstr += "#undef OB_VAR\n"
	retstr += "#define OB_VAR(x) x\n"
	return retstr
}

func SplitVar(sline string) []string {
	var vars []string
	vars = strings.Split(sline, ",")
	for i, v := range vars {
		vars[i] = strings.Trim(v, " \t")
	}
	return vars
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

func ExpandNTabs(ntabs int) (retstr string) {
	retstr = ""
	for i := 0; i < ntabs; i++ {
		retstr += "\t"
	}
	return retstr
}

func ObCodeTransition(line string, codename string, prefix string, ntabs int) (retstr string) {
	var tmpvars []string
	var inputvars []string
	var times int
	var nvars, ci, cj int
	inputvars = SplitVar(codename)

	retstr += ExpandNTabs(ntabs)
	retstr += "do {\n"
	if len(inputvars) < 1 {
		retstr += "}while(0);\n"
		return retstr
	}
	nvars = len(inputvars)
	tmpvars = make([]string, nvars)
	for i, _ := range inputvars {
		tmpvars[i] = fmt.Sprintf("__%s_x%d", prefix, i)
		retstr += ExpandNTabs(ntabs + 1)
		retstr += fmt.Sprintf("int %s;\n", tmpvars[i])
	}

	for i, v := range inputvars {
		retstr += ExpandNTabs(ntabs + 1)
		retstr += fmt.Sprintf("%s = (int)%s;\n", tmpvars[i], v)
	}

	times = int(rand.Float32()*200)%100 + 20
	for i := 0; i < times; i++ {
		ci = int(rand.Float32()*100) % nvars
		cj = int(rand.Float32()*100) % nvars

		retstr += ExpandNTabs(ntabs + 1)
		retstr += fmt.Sprintf("%s = %s;\n", tmpvars[ci], tmpvars[cj])
	}

	retstr += ExpandNTabs(ntabs)
	retstr += "}while(0);\n"
	return retstr

}

func ReadWriteFile(fname string, wfname string, prefix string) (repl int, e error) {
	var willwrite int
	repl = 0
	rf, e := os.Open(fname)
	if e != nil {
		Debug("open %s error %v", fname, e)
		return 0, e
	}

	defer rf.Close()
	wf, e := os.Create(wfname)
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
	obvar_reg, e := regexp.Compile(`OB_VAR\(([^\)]+)\)`)
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
		willwrite = 1
		if obfunc_reg.Match(line) {
			r := obfunc_reg.FindStringSubmatch(string(line))
			if !define_reg.Match(line) {
				retstr := ObFuncTransition(string(line), r[2], prefix)
				wf.WriteString(retstr)
				repl++
			}
		} else if obvar_reg.Match(line) && !define_reg.Match(line) {
			r := obvar_reg.FindStringSubmatch(string(line))
			retstr := ObVarTransition(string(line), r[1], prefix)
			willwrite = 0
			wf.WriteString(retstr)
			repl++

		} else if obcode_reg.Match(line) {
			r := obcode_reg.FindStringSubmatch(string(line))
			if !define_reg.Match(line) {
				ntabs := CountTabs(string(line))
				retstr := ObCodeTransition(string(line), r[1], prefix, ntabs)
				wf.WriteString(retstr)
				repl++
			}
		}
		linenum++
		if willwrite > 0 {
			wf.WriteString(string(line) + "\n")
		}

	}
	return repl, nil
}
