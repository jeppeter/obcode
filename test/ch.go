package main

import (
	"fmt"
	"runtime"
	"time"
)

func Getstring(ch chan string, i int) {

	for {
		select {
		case s := <-ch:
			fmt.Printf("get string %s in %d\n", s, i)
		case <-time.After(time.Millisecond * 30):
			fmt.Printf("%d over\n", i)
			return
		}
	}

}

func main() {
	var schan chan string
	schan = make(chan string, 4)

	for i := 0; i < runtime.NumCPU(); i++ {
		go Getstring(schan, i)
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		s := fmt.Sprintf("format %d", i)
		schan <- s
	}
	time.Sleep(time.Millisecond * 300)
	close(schan)
	return
}
