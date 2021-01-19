package diff

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

type Delta struct {
	Offset uint
	Delete uint
	Insert []byte
}

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

func Compact(deltas []*Delta) (results []*Delta) {
	// Combine deltas with same offset or consecutive deletes
	for i := 0; i < len(deltas); {
		first := deltas[i]
		j := i + 1
		for j < len(deltas) {
			next := deltas[j]
			if first.Offset != next.Offset && (len(first.Insert) > 0 || len(next.Insert) > 0 || first.Offset+first.Delete != next.Offset) {
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

func Cost(deltas []*Delta) (cost uint) {
	// Sum cost of deltas
	for _, d := range deltas {
		cost += 1                       // Every delta has a cost of 1
		cost += d.Delete                // Count every removed byte
		cost += uint(len(d.Insert) * 8) // Count every added byte
	}
	return
}

func Diff(a, b []byte) (deltas []*Delta) {
	deltas = diff(a, b, 0, 0)
	// Rebase deltas into sequence
	var change uint
	for _, d := range deltas {
		d.Offset += change
		change -= d.Delete
		change += uint(len(d.Insert))
	}
	return
}

func diff(a, b []byte, x, y uint) []*Delta {
	var deltas []*Delta
	alen := uint(len(a))
	blen := uint(len(b))
	if x >= alen {
		// All remaining in B must be inserted
		if y < blen {
			deltas = append(deltas, &Delta{
				Offset: x,
				Insert: b[y:],
			})
		}
	} else if y >= blen {
		// All remaining in A must be deleted
		if x < alen {
			deltas = append(deltas, &Delta{
				Offset: x,
				Delete: alen - x,
			})
		}
	} else if a[x] == b[y] {
		// Match
		deltas = diff(a, b, x+1, y+1)
	} else {
		// Try Delete
		ddeltas := diff(a, b, x+1, y)
		dcost := Cost(ddeltas)
		// Try Insert
		ideltas := diff(a, b, x, y+1)
		icost := Cost(ideltas)
		// Compare results
		if dcost <= icost {
			deltas = append(deltas, &Delta{
				Offset: x,
				Delete: 1,
			})
			deltas = append(deltas, ddeltas...)
		} else {
			deltas = append(deltas, &Delta{
				Offset: x,
				Insert: []byte{b[y]},
			})
			deltas = append(deltas, ideltas...)
		}
	}
	return Compact(deltas)
}
