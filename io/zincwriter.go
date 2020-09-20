package io

import (
	"bufio"
	"strings"

	"gitlab.com/NeedleInAJayStack/haystack"
)

// ZincWriter handles writing zinc strings that contain nested grids.
type ZincWriter struct {
	out       *bufio.Writer
	gridDepth int
}

func NewZincWriter(out *bufio.Writer) ZincWriter {
	return ZincWriter{
		out:       out,
		gridDepth: 0,
	}
}

// func GridToZincString(grid haystack.Grid, version int) string {
// 	builder := new(strings.Builder)
// 	writer := NewZincWriter(bufio.NewWriter(builder))
// 	writer.WriteVal(val)
// 	writer.out.Flush()
// 	return builder.String()
// }

func ValToZincString(val haystack.Val) string {
	builder := new(strings.Builder)
	writer := NewZincWriter(bufio.NewWriter(builder))
	writer.WriteVal(val)
	writer.out.Flush()
	return builder.String()
}

func (writer *ZincWriter) WriteVal(val haystack.Val) *ZincWriter {
	switch val := val.(type) {
	case haystack.Grid:
		if writer.gridDepth > 0 {
			writer.gridDepth++
			writer.out.WriteString("<<\n")
			val.WriteZincTo(writer.out, writer.gridDepth)
			writer.out.WriteString(">>")
			writer.gridDepth--
		} else {
			writer.gridDepth++
			val.WriteZincTo(writer.out, 0)
		}
	case haystack.List:
		val.WriteZincTo(writer.out)
	case haystack.Dict:
		val.WriteZincTo(writer.out, true)
	case haystack.Str:
		val.WriteZincTo(writer.out)
	case haystack.Uri:
		val.WriteZincTo(writer.out)
	default:
		writer.out.WriteString(val.ToZinc())
	}
	return writer
}

func (writer *ZincWriter) Flush() {
	writer.out.Flush()
}
