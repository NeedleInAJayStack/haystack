package haystack

type List struct {
	items []Val
}

func (list *List) toZinc() string {
	result := "["
	for idx, val := range list.items {
		if idx != 0 {
			result = result + ", "
		}
		result = result + val.toZinc()
	}
	result = result + "]"
	return result
}
