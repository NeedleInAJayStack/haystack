package client

import (
	"gitlab.com/NeedleInAJayStack/haystack"
)

type CallError struct {
	grid haystack.Grid
}

func NewCallError(grid haystack.Grid) CallError {
	return CallError{grid: grid}
}

func (err CallError) Error() string {
	dis := err.grid.Meta().Get("dis")
	switch val := dis.(type) {
	case haystack.Str:
		return val.String()
	default:
		return "Server side error"
	}
}
