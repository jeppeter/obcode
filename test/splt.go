package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

func Debug(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stdout, s)
	return len(s)
}

func SplitVar(sline string) []string {
	var vars []string
	vars = strings.Split(sline, ",")
	for i, v := range vars {
		vars[i] = strings.Trim(v, " \t")
	}
	return vars
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

	Debug("inputvars %s", inputvars)
	retstr += ExpandNTabs(ntabs)
	retstr += "do {\n"
	if len(inputvars) < 1 {
		retstr += "}while(0)\n"
		return retstr
	}
	nvars = len(inputvars)
	Debug("nvars %d", nvars)
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
	Debug("times %d", times)
	for i := 0; i < times; i++ {
		ci = int(rand.Float32()*100) % nvars
		cj = int(rand.Float32()*100) % nvars

		retstr += ExpandNTabs(ntabs + 1)
		retstr += fmt.Sprintf("%s = %s;\n", tmpvars[ci], tmpvars[cj])
	}

	retstr += ExpandNTabs(ntabs)
	retstr += "}while(0)\n"
	Debug("retstr %s", retstr)
	return retstr

}

func main() {
	var retstr string
	rand.Seed(int64(time.Now().Nanosecond()))
	retstr = ObCodeTransition("", os.Args[1], "prefix", 3)
	fmt.Sprintf("%s", retstr)
}
