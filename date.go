package haystack

import "fmt"

type Date struct {
	year  int
	month int
	day   int
}

// Convert object to zinc. Format is YYYY-MM-DD
func (date *Date) toZinc() string {
	return date.encode()
}

func (date *Date) encode() string {
	result := ""
	result = result + fmt.Sprintf("%d", date.year) + "-"
	if date.month < 10 {
		result = result + "0"
	}
	result = result + fmt.Sprintf("%d", date.month) + "-"
	if date.day < 10 {
		result = result + "0"
	}
	result = result + fmt.Sprintf("%d", date.day)
	return result
}
