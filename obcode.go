package main

import (
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strconv"
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

/********************************************
*    ch for the string to handle
*    done for the over down it
********************************************/
func ObcodeRoutine(srcdir string, dstdir string, ch chan string, cpch chan string, done chan int, over chan int, prefix string) (cnt int, e error) {
	var fname string

	cnt = 0
	e = nil
	for {
		select {
		case fname = <-ch:
			_, err := Obcode(srcdir, dstdir, fname, prefix)
			if err != nil {
				Error("Obcode<%s%s%s>  in %s error %v", srcdir, string(os.PathSeparator), fname, prefix, err)
			} else {
				cnt++
			}

		case fname = <-cpch:
			srcfile := srcdir + string(os.PathSeparator) + fname
			dstfile := dstdir + string(os.PathSeparator) + fname
			err := CopyFileContents(srcfile, dstfile)
			if err != nil {
				Error("Copyfile<%s%s%s> in %s error %v", srcdir, string(os.PathSeparator), fname, prefix, err)
			}
		case <-done:
			goto out_chan

		}
	}
out_chan:
	over <- 1
	return cnt, e
}

func MainDispatch(srcdir string, dstdir string, partfile string, ch chan string, cpch chan string, filterlist *list.List) error {
	var sd, dd, curs, curd, nextpart string
	var doing int

	if len(partfile) > 0 {
		sd = srcdir + string(os.PathSeparator) + partfile
		dd = dstdir + string(os.PathSeparator) + partfile
	} else {
		sd = srcdir
		dd = dstdir
	}

	files, e := ioutil.ReadDir(sd)
	if e != nil {
		Error("can not readdir %s <%v>", sd, e)
		return e
	}
	for _, f := range files {
		if f.Mode().IsDir() {
			curd = dd + string(os.PathSeparator) + f.Name()
			curs = sd + string(os.PathSeparator) + f.Name()
			fileinfo, err := os.Lstat(curs)
			if err != nil {
				Error("can not Lstat %s <%v>", curs, err)
				return err
			}
			err = os.MkdirAll(curd, fileinfo.Mode().Perm())
			if err != nil {
				Error("can not mkdir %s <%v>", curd, err)
				return err
			}
			if len(partfile) > 0 {
				nextpart = partfile + string(os.PathSeparator) + f.Name()
			} else {
				nextpart = f.Name()
			}
			err = MainDispatch(srcdir, dstdir, nextpart, ch, cpch, filterlist)
			if err != nil {
				return err
			}
		} else if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			curs = sd + string(os.PathSeparator) + f.Name()
			curd = dd + string(os.PathSeparator) + f.Name()
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
				var s string
				if reflect.ValueOf(e.Value).Kind() == reflect.String {
					s = reflect.ValueOf(e.Value).String()
					if strings.HasSuffix(f.Name(), s) {
						doing = 1
						break
					}
				}
			}

			if len(partfile) > 0 {
				nextpart = partfile + string(os.PathSeparator) + f.Name()
			} else {
				nextpart = f.Name()
			}

			if doing != 0 {
				Debug("%s%s%s for obfucscated", srcdir, string(os.PathSeparator), nextpart)
				ch <- nextpart
			} else {
				Debug("%s%s%s for copy", srcdir, string(os.PathSeparator), nextpart)
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

func NewThrArgs(srcdir string, dstdir string, prefix string, ch chan string) *ThrArgs {
	args := &ThrArgs{}
	args.srcdir = srcdir
	args.dstdir = dstdir
	args.ch = ch
	args.prefix = prefix
	args.done = make(chan int)
	args.over = make(chan int)
	return args
}

func Usage(ec int, format string, a ...interface{}) {
	f := os.Stderr
	if ec == 0 {
		f = os.Stdout
	}

	if format != "" {
		fmt.Fprintf(f, format, a)
		fmt.Fprintf(f, "\n")
	}

	fmt.Fprintf(f, "obcode [OPTIONS] srddir dstdir\n")
	fmt.Fprintf(f, "\t-h|--help                  :do display this help information\n")
	fmt.Fprintf(f, "\t-f|--filter suffix         :to filter the file\n")
	fmt.Fprintf(f, "\t-p|--prefix prefix         :to set prefix default is \"prefix\"\n")
	fmt.Fprintf(f, "\t-n|--numcpu num            :to specify the number of routines default is num cpu\n")
	os.Exit(ec)
}

var filterlist list.List
var gprefix string
var numroutine int
var gsrcdir, gdstdir string

func ParseArgs() {
	var i int
	var err error
	for i = 1; i < len(os.Args); i++ {
		if os.Args[i] == "-h" || os.Args[i] == "--help" {
			Usage(0, "")
		} else if os.Args[i] == "-f" || os.Args[i] == "--filter" {
			if (i + 1) >= len(os.Args) {
				Usage(3, "%s need args", os.Args[i])
			}
			i++
			filterlist.PushBack(os.Args[i])
		} else if os.Args[i] == "-p" || os.Args[i] == "--prefix" {
			if (i + 1) >= len(os.Args) {
				Usage(3, "%s need args", os.Args[i])
			}
			i++
			gprefix = os.Args[i]
		} else if os.Args[i] == "-n" || os.Args[i] == "--numcpu" {
			if (i + 1) >= len(os.Args) {
				Usage(3, "%s need args", os.Args[i])
			}
			i++
			numroutine, err = strconv.Atoi(os.Args[i])
			if err != nil {
				Usage(3, "%s must numeric", os.Args[i])
			}
		} else {
			break
		}
	}

	if (i + 2) > len(os.Args) {
		Usage(3, "need srcdir dstdir")
	}

	if numroutine <= 0 {
		Usage(3, "num of routine must >0 ")
	}

	gsrcdir = os.Args[i]
	i++
	gdstdir = os.Args[i]
}

func main() {
	var argsarr []*ThrArgs

	gprefix = "prefix"
	numroutine = runtime.NumCPU()
	filterlist.PushBack(".c")
	filterlist.PushBack(".h")
	ParseArgs()
	argsarr = make([]*ThrArgs, numroutine)
	cpch := make(chan string)
	obch := make(chan string)

	_, e := os.Lstat(gdstdir)
	if e != nil {
		fi, _ := os.Lstat(gsrcdir)
		os.MkdirAll(gdstdir, fi.Mode().Perm())
	}

	for i := 0; i < numroutine; i++ {
		s := fmt.Sprintf("%s_%d", gprefix, i)
		args := NewThrArgs(gsrcdir, gdstdir, s, obch)
		argsarr[i] = args
		go ObcodeRoutine(gsrcdir, gdstdir, obch, cpch, args.done, args.over, s)
	}

	err := MainDispatch(gsrcdir, gdstdir, "", obch, cpch, &filterlist)

	for i := 0; i < numroutine; i++ {
		argsarr[i].done <- 1
	}

	for i := 0; i < numroutine; i++ {
		<-argsarr[i].over
	}

	if err != nil {
		os.Exit(3)
	}
	os.Exit(0)
}
