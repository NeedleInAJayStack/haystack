package haystack

import "fmt"

type Str struct {
	val string
}

func NewStr(val string) Str {
	return Str{val: val}
}

// ToZinc representes the object as a double-quoted string, with back-slash escapes
func (str Str) ToZinc() string {
	result := "\""

	for i := 0; i < len(str.val); i++ {
		char := str.val[i]
		if char < ' ' || char == '"' || char == '\\' {
			result = result + "\\"
			switch char {
			case '\n':
				result = result + string('n')
			case '\r':
				result = result + string('r')
			case '\t':
				result = result + string('t')
			case '"':
				result = result + string('"')
			case '\\':
				result = result + string('\\')
			default:
				result = result + "u00"
				if char <= 0xf {
					result = result + string('0')
				}
				result = result + fmt.Sprintf("%x", char)
			}
		} else {
			result = result + string(char)
		}
	}
	result = result + "\""

	return result
}
