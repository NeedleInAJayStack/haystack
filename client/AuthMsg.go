package client

import "strings"

// This models a message in the Haystack authorization header format.
// They follow the form: "[scheme] <name1>=<val1>, <name2>=<val2>, ..."
type authMsg struct {
	scheme string
	attrs  map[string]string
}

func authMsgFromString(str string) authMsg {
	attrs := make(map[string]string)
	attributeStrs := strings.Split(str, ",")
	scheme := ""

	// The first one MAY include the scheme but not necessarily. Handle both situations
	firstAttr := attributeStrs[0]
	if strings.Contains(attributeStrs[0], " ") {
		schemeSplit := strings.Split(firstAttr, " ")
		scheme = strings.TrimSpace(schemeSplit[0])
		attributeStrs[0] = schemeSplit[1]
	}

	for _, attributeStr := range attributeStrs {
		attributeSplit := strings.Split(attributeStr, "=")
		name := strings.TrimSpace(attributeSplit[0])
		val := strings.TrimSpace(attributeSplit[1])
		attrs[name] = val
	}

	return authMsg{
		scheme: scheme,
		attrs:  attrs,
	}
}

func (authMsg *authMsg) get(attrName string) string {
	return authMsg.attrs[attrName]
}

func (authMsg *authMsg) toString() string {
	builder := new(strings.Builder)
	if authMsg.scheme != "" {
		builder.WriteString(strings.ToUpper(authMsg.scheme))
		builder.WriteRune(' ')
	}
	firstVal := true
	for name, val := range authMsg.attrs {
		if firstVal {
			firstVal = false
		} else {
			builder.WriteString(", ")
		}
		builder.WriteString(name)
		builder.WriteRune('=')
		builder.WriteString(val)
	}
	return builder.String()
}
