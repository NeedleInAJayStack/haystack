package haystack

import (
	"bufio"
	"encoding/json"
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

// MarshalJSON representes the object as: "u:<val>"
func (uri Uri) MarshalJSON() ([]byte, error) {
	return json.Marshal("u:" + uri.val)
}

// MarshalHayson representes the object as: "{\"_kind\":\"uri\",\"val\":\"<val>\"}"
func (uri Uri) MarshalHayson() ([]byte, error) {
	return []byte("{\"_kind\":\"uri\",\"val\":\"" + uri.val + "\"}"), nil
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
