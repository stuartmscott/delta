package main

import (
	"github.com/stuartmscott/diff"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: diff a b")
	}
	a := []byte(os.Args[1])
	b := []byte(os.Args[2])
	buffer := a
	for i, d := range diff.Diff(a, b) {
		buffer = diff.Apply(buffer, d)
		log.Println(i, "Offset:", d.Offset, "Delete:", d.Delete, "Insert:", string(d.Insert), "Result:", string(buffer))
	}
}
