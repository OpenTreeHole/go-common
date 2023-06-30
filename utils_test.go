package common

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"testing"
)

func TestMax(t *testing.T) {
	var a, b = 1, 2
	assert.EqualValues(t, 2, Max(a, b))

	var strA, strB = "a", "b"
	assert.EqualValues(t, "b", Max(strA, strB))
}

func TestMin(t *testing.T) {
	var a, b = 1, 2
	assert.EqualValues(t, 1, Min(a, b))

	var strA, strB = "a", "b"
	assert.EqualValues(t, "a", Min(strA, strB))
}

func TestKeys(t *testing.T) {
	var m = map[string]int{"a": 1, "b": 2}
	var keys = Keys(m)
	slices.Sort(keys)
	assert.EqualValues(t, []string{"a", "b"}, keys)
}
