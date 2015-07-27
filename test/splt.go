package main

import (
	"fmt"
	"os"
	"strings"
)

func SplitVar(sline string) []string {
	var vars []string
	vars = strings.Split(sline, ",")
	for i, v := range vars {
		vars[i] = strings.Trim(v, " \t")
	}
	return vars
}
func main() {
	var vars []string

	vars = SplitVar(os.Args[1])
	fmt.Printf("vars %v\n", vars)
}
