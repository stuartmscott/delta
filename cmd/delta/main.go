package main

import (
	"fmt"
	"github.com/stuartmscott/delta"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: delta file1 file2")
		return
	}
	a, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	delta.WriteTo(os.Stdout, delta.Deltas(a, b), a)
}
