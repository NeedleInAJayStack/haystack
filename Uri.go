package haystack

import (
	"bufio"
	"strings"
)

// Uri models a URI tag value.
type Uri struct {
	val string
}

// NewUri creates a new Uri object.
func NewUri(val string) Uri {
	return Uri{val: val}
}

// Type returns the XStr object type
func (uri Uri) String() string {
	return uri.val
}

// ToJSON representes the object as: "u:<val>"
func (uri Uri) ToJSON() string {
	return "u:" + uri.val
}

// ToZinc representes the object as: "`<val>`" with escaped backticks
func (uri Uri) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	uri.WriteZincTo(out)
	out.Flush()
	return builder.String()
}

// WriteZincTo appends the Writer with the URI zinc representation
func (uri Uri) WriteZincTo(buf *bufio.Writer) {
	buf.WriteRune('`')
	for i := 0; i < len(uri.val); i++ {
		char := uri.val[i]
		if char == '`' {
			buf.WriteString("\\`")
		} else if char > ' ' { // URIs cannot contain characters < ' ', so just ignore them.
			buf.WriteByte(char)
		}
	}
	buf.WriteRune('`')
}
