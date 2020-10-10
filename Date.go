package haystack

import (
	"fmt"
	"strconv"
	"strings"
)

// Date models a date (day in year) tag value.
type Date struct {
	year  int
	month int
	day   int
}

// NewDate creates a new Date object. The values are not validated for correctness.
func NewDate(year int, month int, day int) Date {
	return Date{
		year:  year,
		month: month,
		day:   day,
	}
}

// NewDateFromString creates a Date object from a string in the format: "YYYY-MM-DD"
func NewDateFromString(str string) (Date, error) {
	parts := strings.Split(str, "-")

	year, yearErr := strconv.Atoi(parts[0])
	if yearErr != nil {
		return Date{}, yearErr
	}
	month, monthErr := strconv.Atoi(parts[1])
	if monthErr != nil {
		return Date{}, monthErr
	}
	day, dayErr := strconv.Atoi(parts[2])
	if dayErr != nil {
		return Date{}, dayErr
	}

	return Date{
		year:  year,
		month: month,
		day:   day,
	}, nil
}

// Year returns the years of the object.
func (date Date) Year() int {
	return date.year
}

// Month returns the numerical month of the object.
func (date Date) Month() int {
	return date.month
}

// Day returns the day-of-month of the object.
func (date Date) Day() int {
	return date.day
}

// ToZinc representes the object as: "YYYY-MM-DD"
func (date Date) ToZinc() string {
	return date.encode()
}

// ToJSON representes the object as: "d:YYYY-MM-DD"
func (date Date) ToJSON() string {
	return "d:" + date.encode()
}

func (date Date) encode() string {
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

func (date1 *Date) equals(date2 *Date) bool {
	return date1.year == date2.year &&
		date1.month == date2.month &&
		date1.day == date2.day
}
