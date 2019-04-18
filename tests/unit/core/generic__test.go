package coreunit__test

import (
	"testing"

	"github.com/bojodimitrov/byfiri/util"
)

var containstests = []struct {
	name      string
	container []int
	element   int
	result    bool
}{
	{
		name:      "positive",
		container: []int{0, 1, 2},
		element:   1,
		result:    true,
	},
	{
		name:      "negative",
		container: []int{0, 1, 2},
		element:   4,
		result:    false,
	},
}

func TestContains(t *testing.T) {
	for _, ct := range containstests {
		t.Run(ct.name, func(t *testing.T) {
			result := util.Contains(ct.container, ct.element)
			if result != ct.result {
				t.Errorf("got %t, want %t", result, ct.result)
			}
		})
	}
}

var mintests = []struct {
	name   string
	a      int
	b      int
	result int
}{
	{
		name:   "positive",
		a:      2,
		b:      4,
		result: 2,
	},
	{
		name:   "negative",
		a:      -2,
		b:      -4,
		result: -4,
	},
	{
		name:   "mixed",
		a:      2,
		b:      -4,
		result: -4,
	},
}

func TestMin(t *testing.T) {
	for _, mt := range mintests {
		t.Run(mt.name, func(t *testing.T) {
			result := util.Min(mt.a, mt.b)
			if result != mt.result {
				t.Errorf("got %q, want %q", result, mt.result)
			}
		})
	}
}
