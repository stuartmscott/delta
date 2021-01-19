package main

import (
	"fmt"
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
		output := []interface{}{i + 1, "Offset:", d.Offset}
		if d.Delete > 0 {
			output = append(output, "Delete:", d.Delete)
		}
		if len(d.Insert) > 0 {
			output = append(output, "Insert:", fmt.Sprintf("'%s'", d.Insert))
		}
		output = append(output, "Result:", string(buffer))
		log.Println(output...)
	}
}
