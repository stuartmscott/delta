package delta_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stuartmscott/delta"
	"testing"
)

func TestApply(t *testing.T) {
	for name, tt := range map[string]struct {
		given    string
		deltas   []*delta.Delta
		expected string
	}{
		"empty": {},
		"equal": {
			given:    "foobar",
			expected: "foobar",
		},
		"insert_prefix": {
			given: "bar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("foo"),
				},
			},
			expected: "foobar",
		},
		"insert_infix": {
			given: "foar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Offset: 2,
					Insert: []byte("ob"),
				},
			},
			expected: "foobar",
		},
		"insert_suffix": {
			given: "foo",
			deltas: []*delta.Delta{
				&delta.Delta{
					Offset: 3,
					Insert: []byte("bar"),
				},
			},
			expected: "foobar",
		},
		"delete_prefix": {
			given: "foobar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Delete: 3,
				},
			},
			expected: "bar",
		},
		"delete_infix": {
			given: "foobar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Offset: 2,
					Delete: 2,
				},
			},
			expected: "foar",
		},
		"delete_suffix": {
			given: "foobar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Offset: 3,
					Delete: 3,
				},
			},
			expected: "foo",
		},
		"swap": {
			given: "foobar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("bar"),
				},
				&delta.Delta{
					Offset: 6,
					Delete: 3,
				},
			},
			expected: "barfoo",
		},
		"delete_vowels": {
			given: "foobar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Offset: 1,
					Delete: 2,
				},
				&delta.Delta{
					Offset: 2,
					Delete: 1,
				},
			},
			expected: "fbr",
		},
		"delete_consonants": {
			given: "foobar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
				},
				&delta.Delta{
					Offset: 2,
					Delete: 1,
				},
				&delta.Delta{
					Offset: 3,
					Delete: 1,
				},
			},
			expected: "ooa",
		},
		"insert_vowels": {
			given: "fbr",
			deltas: []*delta.Delta{
				&delta.Delta{
					Offset: 1,
					Insert: []byte("oo"),
				},
				&delta.Delta{
					Offset: 4,
					Insert: []byte("a"),
				},
			},
			expected: "foobar",
		},
		"insert_consonants": {
			given: "ooa",
			deltas: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("f"),
				},
				&delta.Delta{
					Offset: 3,
					Insert: []byte("b"),
				},
				&delta.Delta{
					Offset: 5,
					Insert: []byte("r"),
				},
			},
			expected: "foobar",
		},
		"replace": {
			given: "foo",
			deltas: []*delta.Delta{
				&delta.Delta{
					Delete: 3,
					Insert: []byte("bar"),
				},
			},
			expected: "bar",
		},
		"reverse": {
			given: "foobar",
			deltas: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
					Insert: []byte("rab"),
				},
				&delta.Delta{
					Offset: 5,
					Delete: 3,
					Insert: []byte("f"),
				},
			},
			expected: "raboof",
		},
	} {
		t.Run(name, func(t *testing.T) {
			buffer := []byte(tt.given)
			for _, d := range tt.deltas {
				buffer = delta.Apply(buffer, d)
			}
			assert.Equal(t, tt.expected, string(buffer))
		})
	}
}

func TestCompact(t *testing.T) {
	for name, tt := range map[string]struct {
		deltas, expected []*delta.Delta
	}{
		"empty": {},
		"single": {
			deltas: []*delta.Delta{
				&delta.Delta{},
			},
			expected: []*delta.Delta{
				&delta.Delta{},
			},
		},
		"consecutive": {
			deltas: []*delta.Delta{
				&delta.Delta{},
				&delta.Delta{
					Offset: 1,
				},
			},
			expected: []*delta.Delta{
				&delta.Delta{},
				&delta.Delta{
					Offset: 1,
				},
			},
		},
		"delete_delete": {
			deltas: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
				},
				&delta.Delta{
					Offset: 1,
					Delete: 1,
				},
			},
			expected: []*delta.Delta{
				&delta.Delta{
					Delete: 2,
				},
			},
		},
		"insert_insert": {
			deltas: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("a"),
				},
				&delta.Delta{
					Insert: []byte("b"),
				},
			},
			expected: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("ab"),
				},
			},
		},
		"delete_insert": {
			deltas: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
				},
				&delta.Delta{
					Insert: []byte("a"),
				},
			},
			expected: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
					Insert: []byte("a"),
				},
			},
		},
		"insert_delete": {
			deltas: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("a"),
				},
				&delta.Delta{
					Delete: 1,
				},
			},
			expected: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
					Insert: []byte("a"),
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := delta.Compact(tt.deltas)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCost(t *testing.T) {
	for name, tt := range map[string]struct {
		deltas   []*delta.Delta
		expected uint
	}{
		"empty": {
			expected: 0,
		},
		"single": {
			deltas: []*delta.Delta{
				&delta.Delta{},
			},
			expected: 1,
		},
		"delete": {
			deltas: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
				},
			},
			expected: 2,
		},
		"insert": {
			deltas: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("a"),
				},
			},
			expected: 9,
		},
		"replace": {
			deltas: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
					Insert: []byte("a"),
				},
			},
			expected: 10,
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := delta.Cost(tt.deltas)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestDeltas(t *testing.T) {
	for name, tt := range map[string]struct {
		a, b     string
		expected []*delta.Delta
	}{
		"empty": {},
		"equal": {
			a: "foobar",
			b: "foobar",
		},
		"insert_prefix": {
			a: "bar",
			b: "foobar",
			expected: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("foo"),
				},
			},
		},
		"insert_infix": {
			a: "foar",
			b: "foobar",
			expected: []*delta.Delta{
				&delta.Delta{
					Offset: 2,
					Insert: []byte("ob"),
				},
			},
		},
		"insert_suffix": {
			a: "foo",
			b: "foobar",
			expected: []*delta.Delta{
				&delta.Delta{
					Offset: 3,
					Insert: []byte("bar"),
				},
			},
		},
		"delete_prefix": {
			a: "foobar",
			b: "bar",
			expected: []*delta.Delta{
				&delta.Delta{
					Delete: 3,
				},
			},
		},
		"delete_infix": {
			a: "foobar",
			b: "foar",
			expected: []*delta.Delta{
				&delta.Delta{
					Offset: 2,
					Delete: 2,
				},
			},
		},
		"delete_suffix": {
			a: "foobar",
			b: "foo",
			expected: []*delta.Delta{
				&delta.Delta{
					Offset: 3,
					Delete: 3,
				},
			},
		},
		"swap": {
			a: "foobar",
			b: "barfoo",
			expected: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("bar"),
				},
				&delta.Delta{
					Offset: 6,
					Delete: 3,
				},
			},
		},
		"delete_vowels": {
			a: "foobar",
			b: "fbr",
			expected: []*delta.Delta{
				&delta.Delta{
					Offset: 1,
					Delete: 2,
				},
				&delta.Delta{
					Offset: 2,
					Delete: 1,
				},
			},
		},
		"delete_consonants": {
			a: "foobar",
			b: "ooa",
			expected: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
				},
				&delta.Delta{
					Offset: 2,
					Delete: 1,
				},
				&delta.Delta{
					Offset: 3,
					Delete: 1,
				},
			},
		},
		"insert_vowels": {
			a: "fbr",
			b: "foobar",
			expected: []*delta.Delta{
				&delta.Delta{
					Offset: 1,
					Insert: []byte("oo"),
				},
				&delta.Delta{
					Offset: 4,
					Insert: []byte("a"),
				},
			},
		},
		"insert_consonants": {
			a: "ooa",
			b: "foobar",
			expected: []*delta.Delta{
				&delta.Delta{
					Insert: []byte("f"),
				},
				&delta.Delta{
					Offset: 3,
					Insert: []byte("b"),
				},
				&delta.Delta{
					Offset: 5,
					Insert: []byte("r"),
				},
			},
		},
		"replace": {
			a: "foo",
			b: "bar",
			expected: []*delta.Delta{
				&delta.Delta{
					Delete: 3,
					Insert: []byte("bar"),
				},
			},
		},
		"reverse": {
			a: "foobar",
			b: "raboof",
			expected: []*delta.Delta{
				&delta.Delta{
					Delete: 1,
					Insert: []byte("rab"),
				},
				&delta.Delta{
					Offset: 5,
					Delete: 3,
					Insert: []byte("f"),
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := delta.Deltas([]byte(tt.a), []byte(tt.b))
			assert.Equal(t, tt.expected, actual)
		})
	}
}
