package io

import (
	"testing"

	"github.com/NeedleInAJayStack/haystack"
	"github.com/stretchr/testify/assert"
)

func TestGridFromZinc(t *testing.T) {
	actual, err := GridFromZinc(`ver:"2.0"
fooBar33
`)
	if err != nil {
		t.Error(err)
	}

	gb := haystack.NewGridBuilder()
	gb.AddCol(
		"fooBar33",
		map[string]haystack.Val{},
	)
	expected := gb.ToGrid()

	assert.Equal(t, actual, expected)
}
