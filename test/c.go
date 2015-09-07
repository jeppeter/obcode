package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

type ThrArgs struct {
	fnamechan chan string
	ended     chan int
	closed    chan int
	srcdir    string
	dstdir    string
}

func Thread(num int, args *ThrArgs) (cnt int, e error) {
	cnt = 0

	for {
		select {
		case fname := <-args.fnamechan:
			fmt.Printf("<%d:cnt:%d> srcdir<%s> dstdir<%s> file<%s>\n", num, cnt, args.srcdir, args.dstdir, fname)
			if (rand.Int31() % 2) == 1 {
				time.Sleep(time.Millisecond * 1)
			}
			cnt++
		case <-args.ended:
			fmt.Printf("<%d>got ended\n", num)
			goto out
		}
	}
out:
	args.closed <- 1
	return
}

func NewThreadArgs(srcdir string, dstdir string, fnchan chan string, ended chan int) *ThrArgs {
	args := &ThrArgs{}
	args.srcdir = srcdir
	args.dstdir = dstdir
	args.fnamechan = fnchan
	args.ended = ended
	args.closed = make(chan int)
	return args
}

func main() {
	var argsarr []*ThrArgs
	var seed int64
	seed = int64(time.Now().Nanosecond())
	rand.Seed(seed)
	fnamechan := make(chan string)

	argsarr = make([]*ThrArgs, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		endchan := make(chan int)
		argsarr[i] = NewThreadArgs(os.Args[1], os.Args[2], fnamechan, endchan)
		go Thread(i, argsarr[i])
	}

	for i := 0; i < runtime.NumCPU()*10; i++ {
		str := fmt.Sprintf("File%d", i)
		fnamechan <- str
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		argsarr[i].ended <- 1
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		<-argsarr[i].closed
	}

	return
}
