package main

import (
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func Debug(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stdout, s)
	return len(s)
}

func Error(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stderr, s)
	return len(s)
}

func Obcode(srcdir string, dstdir string, fname string, prefix string) (replc int, e error) {
	sfile := srcdir + string(os.PathSeparator) + fname
	dfile := dstdir + string(os.PathSeparator) + fname
	return ReadWriteFile(sfile, dfile, prefix)
}

/********************************************
*    ch for the string to handle
*    done for the over down it
********************************************/
func ObcodeRoutine(srcdir string, dstdir string, ch chan string, done chan int, over chan int, prefix string) (cnt int, e error) {
	var fname string
	cnt = 0
	e = nil
	for {
		select {
		case fname = <-ch:
			repl, err := Obcode(srcdir, dstdir, fname, prefix)
			if err != nil {
				Error("Obcode<%s%c%s>  in %s error %v", srcdir, os.PathSeparator, fname, prefix, err)
				e = err
				goto out_chan
			} else {
				cnt++
			}
		case <-done:
			goto out_chan

		}
	}
out_chan:
	over <- 1
	return cnt, e
}

func CopyFileContents(srcfile string, dstfile string) error {
	in, err := os.Open(srcfile)
	if err != nil {
		Error("could not open %s <%v>", srcfile, err)
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dstfile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Error("could not open %s for write <%v>", dstfile, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		Error("could not copy <%s => %s> <%v>", srcfile, dstfile, err)
		return err
	}

	err = out.Sync()
	if err != nil {
		Error("could not sync %s <%v>", dstfile, err)
		return err
	}
	return nil
}

func CopyFileRoutine(srcdir string, dstdir string, ch chan string, done chan int, over chan int) (cnt int, e error) {
	var fname string

	cnt = 0
	e = nil
	for {
		select {
		case fname = <-ch:
			srcfile := srcdir + os.PathSeparator + fname
			dstfile := dstdir + os.PathSeparator + fname
			err := CopyFileContents(srcfile, dstfile)
			if err != nil {
				e = err
				goto out_chan2
			}
			cnt++
		case <-done:
			goto out_chan2
		}
	}
out_chan2:
	over <- 1
	return cnt, e
}

func MainDispatch(srcdir string, dstdir string, partfile string, ch chan string, cpch chan string, filterlist *list.List) error {
	var sd, dd, curs, curd, nextpart string
	var err error
	var doing int

	if len(partfile) > 0 {
		sd = srcdir + os.PathSeparator + partfile
		dd = dstdir + os.PathSeparator + partfile
	} else {
		sd = srcdir
		dd = dstdir
	}

	files, e := ioutil.ReadDir(sd)
	for i, f := range files {
		if f.Mode().IsDir() {
			curd = dd + os.PathSeparator + f.Name()
			curs = sd + os.PathSeparator + f.Name()
			perm, err := os.Lstat(curs)
			if err != nil {
				Error("can not Lstat %s <%v>", curs, err)
				return err
			}
			err = os.MkdirAll(curd, perm)
			if err != nil {
				Error("can not mkdir %s <%v>", curd, err)
				return err
			}
			if len(partfile) > 0 {
				nextpart = partfile + os.PathSeparator + f.Name()
			} else {
				nextpart = f.Name()
			}
			err = MainDispatch(srcdir, dstdir, nextpart, ch)
			if err != nil {
				return err
			}
		} else if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			curs = sd + os.PathSeparator + f.Name()
			curd = dd + os.PathSeparator + f.Name()
			ln, err := os.Readlink(curs)
			if err != nil {
				Error("could not readlink %s <%v>", curs, err)
				return err
			}

			err = os.Symlink(curd, ln)
			if err != nil {
				Error("could not symlink <%s> to <%s> <%v>", curd, ln, err)
				return err
			}

		} else if f.Mode().IsRegular() {
			doing = 0
			for e := filterlist.Front(); e != nil; e = e.Next() {
				if strings.HasSuffix(f.Name(), e.Value) {
					doing = 1
					break
				}
			}

			if len(partfile) > 0 {
				nextpart = partfile + os.PathSeparator + f.Name()
			} else {
				nextpart = f.Name()
			}

			if doing {
				ch <- nextpart
			} else {
				cpch <- nextpart
			}
		}
	}
	return nil
}

type ThrArgs struct {
	ch     chan string
	done   chan int
	over   chan int
	srcdir string
	dstdir string
	prefix string
}

func NewThrArgs(srcdir string, dstdir string, prefix string, ch chan string) {
	args := &ThrArgs{}
	args.srcdir = srcdir
	args.dstdir = dstdir
	args.ch = ch
	args.prefix = prefix
	args.done = make(chan int)
	args.over = make(chan int)
	return args
}

func main() {
	var argsarr *[]ThrArgs

	if len(os.Args) < 3 {

	}

	argsarr = make(*ThrArgs, runtime.NumCPU()*2)
	cpch := make(chan string)
	obch := make(chan string)

}
