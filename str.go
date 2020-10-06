package haystack

import (
	"bufio"
	"fmt"
	"strings"
)

// Str models a string tag value.
type Str struct {
	val string
}

// NewStr creates a new Str object
func NewStr(val string) Str {
	return Str{val: val}
}

// String returns the object's value directly as a Go string
func (str Str) String() string {
	return str.val
}

// ToZinc representes the object as a double-quoted string, with back-slash escapes
func (str Str) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	str.WriteZincTo(out)
	out.Flush()
	return builder.String()
}

// WriteZincTo writes the object as a double-quoted string, with back-slash escapes to the given writer
func (str Str) WriteZincTo(buf *bufio.Writer) {
	buf.WriteRune('"')

	for i := 0; i < len(str.val); i++ {
		char := str.val[i]
		if char < ' ' || char == '"' || char == '\\' {
			buf.WriteRune('\\')
			switch char {
			case '\n':
				buf.WriteRune('n')
			case '\r':
				buf.WriteRune('r')
			case '\t':
				buf.WriteRune('t')
			case '"':
				buf.WriteRune('"')
			case '\\':
				buf.WriteRune('\\')
			default:
				buf.WriteString("u00")
				if char <= 0xf {
					buf.WriteRune('0')
				}
				buf.WriteString(fmt.Sprintf("%x", char))
			}
		} else {
			buf.WriteByte(char)
		}
	}
	buf.WriteRune('"')
}
