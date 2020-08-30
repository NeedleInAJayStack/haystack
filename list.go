package haystack

type List struct {
	items []Val
}

func (list *List) toZinc() string {
	result := "["
	for i := 0; i < len(list.items); i++ {
		if i > 0 {
			result = result + ","
		}
		item := list.items[i]
		result = result + item.toZinc()
	}
	result = result + "]"
	return result
}
