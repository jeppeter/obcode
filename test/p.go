package main

import (
	"./mypkg"
	"fmt"
)

func main() {
	fmt.Printf("%d + %d = %d\n", 1, 2, mypkg.Add(1, 2))
}
