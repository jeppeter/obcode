package main

import (
	"./stack"
	"fmt"
	"os"
)

func main() {
	var argstk *stack.Stack
	argstk = new(stack.Stack)
	for _, v := range os.Args[1:] {
		argstk.Push(v)
	}

	for {
		v := argstk.Pop()
		if v == nil {
			break
		}
		fmt.Printf("value %s\n", v)
	}
	return
}
