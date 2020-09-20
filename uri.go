package haystack

import (
	"bufio"
	"strings"
)

type Uri struct {
	val string
}

func NewUri(val string) Uri {
	return Uri{val: val}
}

func (uri Uri) String() string {
	return uri.val
}

func (uri Uri) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	uri.WriteZincTo(out)
	out.Flush()
	return builder.String()
}

// ToZinc representes the object as: "`<val>`"
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
