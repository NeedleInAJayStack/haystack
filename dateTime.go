package haystack

type DateTime struct {
	date Date
	time Time
	tz   string
}

// Format is YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz
func (dateTime *DateTime) toZinc() string {
	return dateTime.encode()
}

func (dateTime *DateTime) encode() string {
	result := ""
	result = result + dateTime.date.encode()
	result = result + "T"
	result = result + dateTime.time.encode()

	return result
}
