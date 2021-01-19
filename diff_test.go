package diff_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stuartmscott/diff"
	"testing"
)

func TestApply(t *testing.T) {
	for name, tt := range map[string]struct {
		given    string
		deltas   []*diff.Delta
		expected string
	}{
		"empty": {},
		"equal": {
			given:    "foobar",
			expected: "foobar",
		},
		"insert_prefix": {
			given: "bar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("foo"),
				},
			},
			expected: "foobar",
		},
		"insert_infix": {
			given: "foar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Offset: 2,
					Insert: []byte("ob"),
				},
			},
			expected: "foobar",
		},
		"insert_suffix": {
			given: "foo",
			deltas: []*diff.Delta{
				&diff.Delta{
					Offset: 3,
					Insert: []byte("bar"),
				},
			},
			expected: "foobar",
		},
		"delete_prefix": {
			given: "foobar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Delete: 3,
				},
			},
			expected: "bar",
		},
		"delete_infix": {
			given: "foobar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Offset: 2,
					Delete: 2,
				},
			},
			expected: "foar",
		},
		"delete_suffix": {
			given: "foobar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Offset: 3,
					Delete: 3,
				},
			},
			expected: "foo",
		},
		"swap": {
			given: "foobar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("bar"),
				},
				&diff.Delta{
					Offset: 6,
					Delete: 3,
				},
			},
			expected: "barfoo",
		},
		"delete_vowels": {
			given: "foobar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Offset: 1,
					Delete: 2,
				},
				&diff.Delta{
					Offset: 2,
					Delete: 1,
				},
			},
			expected: "fbr",
		},
		"delete_consonants": {
			given: "foobar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
				},
				&diff.Delta{
					Offset: 2,
					Delete: 1,
				},
				&diff.Delta{
					Offset: 3,
					Delete: 1,
				},
			},
			expected: "ooa",
		},
		"insert_vowels": {
			given: "fbr",
			deltas: []*diff.Delta{
				&diff.Delta{
					Offset: 1,
					Insert: []byte("oo"),
				},
				&diff.Delta{
					Offset: 4,
					Insert: []byte("a"),
				},
			},
			expected: "foobar",
		},
		"insert_consonants": {
			given: "ooa",
			deltas: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("f"),
				},
				&diff.Delta{
					Offset: 3,
					Insert: []byte("b"),
				},
				&diff.Delta{
					Offset: 5,
					Insert: []byte("r"),
				},
			},
			expected: "foobar",
		},
		"replace": {
			given: "foo",
			deltas: []*diff.Delta{
				&diff.Delta{
					Delete: 3,
					Insert: []byte("bar"),
				},
			},
			expected: "bar",
		},
		"reverse": {
			given: "foobar",
			deltas: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
					Insert: []byte("rab"),
				},
				&diff.Delta{
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
				buffer = diff.Apply(buffer, d)
			}
			assert.Equal(t, tt.expected, string(buffer))
		})
	}
}

func TestCompact(t *testing.T) {
	for name, tt := range map[string]struct {
		deltas, expected []*diff.Delta
	}{
		"empty": {},
		"single": {
			deltas: []*diff.Delta{
				&diff.Delta{},
			},
			expected: []*diff.Delta{
				&diff.Delta{},
			},
		},
		"consecutive": {
			deltas: []*diff.Delta{
				&diff.Delta{},
				&diff.Delta{
					Offset: 1,
				},
			},
			expected: []*diff.Delta{
				&diff.Delta{},
				&diff.Delta{
					Offset: 1,
				},
			},
		},
		"delete_delete": {
			deltas: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
				},
				&diff.Delta{
					Offset: 1,
					Delete: 1,
				},
			},
			expected: []*diff.Delta{
				&diff.Delta{
					Delete: 2,
				},
			},
		},
		"insert_insert": {
			deltas: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("a"),
				},
				&diff.Delta{
					Insert: []byte("b"),
				},
			},
			expected: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("ab"),
				},
			},
		},
		"delete_insert": {
			deltas: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
				},
				&diff.Delta{
					Insert: []byte("a"),
				},
			},
			expected: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
					Insert: []byte("a"),
				},
			},
		},
		"insert_delete": {
			deltas: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("a"),
				},
				&diff.Delta{
					Delete: 1,
				},
			},
			expected: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
					Insert: []byte("a"),
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := diff.Compact(tt.deltas)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCost(t *testing.T) {
	for name, tt := range map[string]struct {
		deltas   []*diff.Delta
		expected uint
	}{
		"empty": {
			expected: 0,
		},
		"single": {
			deltas: []*diff.Delta{
				&diff.Delta{},
			},
			expected: 1,
		},
		"delete": {
			deltas: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
				},
			},
			expected: 2,
		},
		"insert": {
			deltas: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("a"),
				},
			},
			expected: 9,
		},
		"replace": {
			deltas: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
					Insert: []byte("a"),
				},
			},
			expected: 10,
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := diff.Cost(tt.deltas)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestDiff(t *testing.T) {
	for name, tt := range map[string]struct {
		a, b     string
		expected []*diff.Delta
	}{
		"empty": {},
		"equal": {
			a: "foobar",
			b: "foobar",
		},
		"insert_prefix": {
			a: "bar",
			b: "foobar",
			expected: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("foo"),
				},
			},
		},
		"insert_infix": {
			a: "foar",
			b: "foobar",
			expected: []*diff.Delta{
				&diff.Delta{
					Offset: 2,
					Insert: []byte("ob"),
				},
			},
		},
		"insert_suffix": {
			a: "foo",
			b: "foobar",
			expected: []*diff.Delta{
				&diff.Delta{
					Offset: 3,
					Insert: []byte("bar"),
				},
			},
		},
		"delete_prefix": {
			a: "foobar",
			b: "bar",
			expected: []*diff.Delta{
				&diff.Delta{
					Delete: 3,
				},
			},
		},
		"delete_infix": {
			a: "foobar",
			b: "foar",
			expected: []*diff.Delta{
				&diff.Delta{
					Offset: 2,
					Delete: 2,
				},
			},
		},
		"delete_suffix": {
			a: "foobar",
			b: "foo",
			expected: []*diff.Delta{
				&diff.Delta{
					Offset: 3,
					Delete: 3,
				},
			},
		},
		"swap": {
			a: "foobar",
			b: "barfoo",
			expected: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("bar"),
				},
				&diff.Delta{
					Offset: 6,
					Delete: 3,
				},
			},
		},
		"delete_vowels": {
			a: "foobar",
			b: "fbr",
			expected: []*diff.Delta{
				&diff.Delta{
					Offset: 1,
					Delete: 2,
				},
				&diff.Delta{
					Offset: 2,
					Delete: 1,
				},
			},
		},
		"delete_consonants": {
			a: "foobar",
			b: "ooa",
			expected: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
				},
				&diff.Delta{
					Offset: 2,
					Delete: 1,
				},
				&diff.Delta{
					Offset: 3,
					Delete: 1,
				},
			},
		},
		"insert_vowels": {
			a: "fbr",
			b: "foobar",
			expected: []*diff.Delta{
				&diff.Delta{
					Offset: 1,
					Insert: []byte("oo"),
				},
				&diff.Delta{
					Offset: 4,
					Insert: []byte("a"),
				},
			},
		},
		"insert_consonants": {
			a: "ooa",
			b: "foobar",
			expected: []*diff.Delta{
				&diff.Delta{
					Insert: []byte("f"),
				},
				&diff.Delta{
					Offset: 3,
					Insert: []byte("b"),
				},
				&diff.Delta{
					Offset: 5,
					Insert: []byte("r"),
				},
			},
		},
		"replace": {
			a: "foo",
			b: "bar",
			expected: []*diff.Delta{
				&diff.Delta{
					Delete: 3,
					Insert: []byte("bar"),
				},
			},
		},
		"reverse": {
			a: "foobar",
			b: "raboof",
			expected: []*diff.Delta{
				&diff.Delta{
					Delete: 1,
					Insert: []byte("rab"),
				},
				&diff.Delta{
					Offset: 5,
					Delete: 3,
					Insert: []byte("f"),
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := diff.Diff([]byte(tt.a), []byte(tt.b))
			assert.Equal(t, tt.expected, actual)
		})
	}
}
