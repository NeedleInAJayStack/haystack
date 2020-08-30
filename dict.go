package haystack

type Dict struct {
	items map[string]Val
}

// Convert object to zinc.
func (dict *Dict) toZinc() string {
	result := "{"
	firstVal := true
	for name, val := range dict.items {
		if firstVal {
			firstVal = false
		} else {
			result = result + ","
		}
		result = result + name + ":" + val.toZinc()
	}
	// for i := 0; i < len(dict.items); i++ {
	// 	if i > 0 {
	// 		result = result + ","
	// 	}
	// 	name, val := dict.items[i]
	// 	result = result + name + ":" + val.toZinc()
	// }
	result = result + "}"
	return result
}
