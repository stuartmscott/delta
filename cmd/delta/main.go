package main

import (
	"flag"
	"fmt"
	"github.com/stuartmscott/delta"
	"io/ioutil"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: delta file1 file2\n")
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		return
	}
	a, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := ioutil.ReadFile(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	delta.WriteTo(os.Stdout, delta.Deltas(a, b), a)
}
