package io

import (
	"errors"

	"github.com/NeedleInAJayStack/haystack"
)

// GridFromZinc parses a Zinc string into a Haystack Grid
func GridFromZinc(zinc string) (haystack.Grid, error) {
	var reader ZincReader
	reader.InitString(zinc)
	val, err := reader.ReadVal()
	if err != nil {
		return haystack.EmptyGrid(), err
	}
	switch val := val.(type) {
	case haystack.Grid:
		return val, nil
	default:
		return (haystack.EmptyGrid()), errors.New("input is not a grid")
	}
}
