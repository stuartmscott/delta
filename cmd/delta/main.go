package main

import (
	"fmt"
	"github.com/stuartmscott/delta"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: delta a b")
	}
	a := []byte(os.Args[1])
	b := []byte(os.Args[2])
	buffer := a
	for i, d := range delta.Deltas(a, b) {
		buffer = delta.Apply(buffer, d)
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
