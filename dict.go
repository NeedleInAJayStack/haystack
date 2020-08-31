package haystack

type Dict struct {
	items map[string]Val
}

func (dict *Dict) isEmpty() bool {
	return len(dict.items) == 0
}

// Format is {name1:val1 name2:val2 ...}. Markers don't require a :val.
func (dict *Dict) toZinc() string {
	result := dict.encode(true)
	return result
}

func (dict *Dict) encode(brackets bool) string {
	result := ""
	if brackets {
		result = result + "{"
	}
	firstVal := true
	for name, val := range dict.items {
		if firstVal {
			firstVal = false
		} else {
			result = result + " "
		}

		_, isMarker := val.(*Marker)
		if isMarker {
			result = result + name
		} else {
			result = result + name + ":" + val.toZinc()
		}
	}
	if brackets {
		result = result + "}"
	}
	return result
}
