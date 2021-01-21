package delta

import (
	"fmt"
	"io"
)

/*
	Eugene W. Myers - An O(ND)Difference Algorithm and Its Variations

	       a
	   0 1 2 3 4 5
	  0*-*-*-*-*-*
	   |\|\|\|\|\|
	  1*-*-*-*-*-*
	   |\|\|\|\|\|
	b 2*-*-*-*-*-*
	   |\|\|\|\|\|
	  3*-*-*-*-*-*
	   |\|\|\|\|\|
	  4*-*-*-*-*-*
	   |\|\|\|\|\|
	  5*-*-*-*-*-*

	* x, y point
	- delete
	| insert
	\ match
*/

// Delta represents an operation (Delete and/or Insert) at an Offset.
type Delta struct {
	Offset uint
	Delete uint
	Insert []byte
}

// Apply modifies the given input with the given Delta.
func Apply(input []byte, delta *Delta) []byte {
	count := 0
	length := uint(len(input))
	output := make([]byte, length-delta.Delete+uint(len(delta.Insert)))
	if delta.Offset <= length {
		count += copy(output, input[:delta.Offset])
	}
	count += copy(output[count:], delta.Insert)
	index := delta.Offset + delta.Delete
	if index < length {
		copy(output[count:], input[index:])
	}
	return output
}

// Compact combines deltas with same offset or consecutive deletes.
func Compact(deltas []*Delta) (results []*Delta) {
	for i := 0; i < len(deltas); {
		first := deltas[i]
		j := i + 1
		for j < len(deltas) {
			next := deltas[j]
			if first.Offset != next.Offset && first.Offset+first.Delete != next.Offset {
				break
			}
			first.Delete += next.Delete
			first.Insert = append(first.Insert, next.Insert...)
			j++
		}
		i = j
		results = append(results, first)
	}
	return
}

// Deltas returns a sequence of deltas that transform the first of the given byte arrays into the second.
func Deltas(a, b []byte) []*Delta {
	n := len(a)
	m := len(b)
	edits := edits(a, b, n, m)
	if edits < 0 {
		// Error determining length of edit script
		return nil
	}
	ds := deltas(a, b, n, m, edits)
	// Compact deltas
	ds = Compact(ds)
	// Rebase deltas into sequence
	var change uint
	for _, d := range ds {
		d.Offset += change
		change -= d.Delete
		change += uint(len(d.Insert))
	}
	return ds
}

func WriteTo(out io.Writer, deltas []*Delta, buffer []byte) {
	for _, d := range deltas {
		fmt.Fprintf(out, "@ %v\n", d.Offset)
		for i := uint(0); i < d.Delete; i++ {
			v := buffer[d.Offset+i]
			fmt.Fprintf(out, "- %3v %s\n", v, string(v))
		}
		buffer = Apply(buffer, d)
		for i := uint(0); i < uint(len(d.Insert)); i++ {
			v := buffer[d.Offset+i]
			fmt.Fprintf(out, "+ %3v %s\n", v, string(v))
		}
	}
}

func edits(a, b []byte, n, m int) int {
	for _, max := range []int{
		minimum(n, m),
		maximum(n, m),
		n + m,
		n * m,
	} {
		v := make(map[int]int, max+2)
		v[1] = 0
		for d := 0; d <= max; d++ {
			start := -(d - 2*maximum(0, d-m))
			end := d - 2*maximum(0, d-n)
			for k := start; k <= end; k = k + 2 {
				x, y := next(a, b, n, m, v, d, k)
				v[k] = x
				if x >= n && y >= m {
					return d
				}
			}
		}
	}
	return -1
}

func deltas(a, b []byte, n, m int, max int) []*Delta {
	v := make(map[int]int, max+2)
	v[1] = 0
	vs := make([]map[int]int, max+1)
	for d := 0; d <= max; d++ {
		start := -(d - 2*maximum(0, d-m))
		end := d - 2*maximum(0, d-n)
		count := (end-start)/2 + 1
		vs[d] = make(map[int]int, count)
		for k := start; k <= end; k = k + 2 {
			x, y := next(a, b, n, m, v, d, k)
			v[k] = x
			vs[d][k] = x
			if x >= n && y >= m {
				return backtrack(a, b, vs, d, n-m)
			}
		}
	}
	return nil
}

func backtrack(a, b []byte, vs []map[int]int, d, k int) []*Delta {
	if d <= 0 {
		return nil
	}
	prev := vs[d-1]
	delta := &Delta{}
	delete, dok := prev[k-1]
	insert, iok := prev[k+1]
	if dok && iok {
		// Select the best
		if delete >= insert {
			iok = false
		} else {
			dok = false
		}
	}
	if dok {
		k = k - 1
		delta.Offset = uint(delete)
		delta.Delete = 1
	} else if iok {
		k = k + 1
		delta.Offset = uint(insert)
		delta.Insert = []byte{b[insert-k]}
	} else {
		// Uh oh
		return nil
	}
	return append(backtrack(a, b, vs, d-1, k), delta)
}

func next(a, b []byte, n, m int, v map[int]int, d, k int) (x int, y int) {
	if k == -d {
		// Left border so choose k above
		x = v[k+1]
	} else if k == d {
		// Top border so choose k below
		x = v[k-1] + 1
	} else {
		above, aok := v[k+1]
		below, bok := v[k-1]
		if aok && bok {
			// Choose best
			if below < above {
				x = above
			} else {
				x = below + 1
			}
		} else if aok {
			// Choose above
			x = above
		} else {
			// Choose below
			x = below + 1
		}
	}
	y = x - k
	for x < n && y < m && a[x] == b[y] {
		x, y = x+1, y+1
	}
	return
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maximum(a, b int) int {
	if a > b {
		return a
	}
	return b
}
