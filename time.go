package haystack

import "fmt"

type Time struct {
	hour int
	min  int
	sec  int
	ms   int // Default value is 0
}

// Format is hh:mm:ss.mmm
func (time *Time) toZinc() string {
	return time.encode()
}

func (time *Time) encode() string {
	result := ""
	if time.hour < 10 {
		result = result + "0"
	}
	result = result + fmt.Sprintf("%d", time.hour) + ":"
	if time.min < 10 {
		result = result + "0"
	}
	result = result + fmt.Sprintf("%d", time.min) + ":"
	if time.sec < 10 {
		result = result + "0"
	}
	result = result + fmt.Sprintf("%d", time.sec)
	if time.ms != 0 {
		result = result + "."
		if time.ms < 10 {
			result = result + "0"
		}
		if time.ms < 100 {
			result = result + "0"
		}
		result = result + fmt.Sprintf("%d", time.ms)
	}
	return result
}
