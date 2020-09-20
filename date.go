package haystack

import (
	"fmt"
	"strconv"
	"strings"
)

type Date struct {
	year  int
	month int
	day   int
}

func dateDef() Date {
	return Date{
		year:  0,
		month: 0,
		day:   0,
	}
}

func NewDate(year int, month int, day int) Date {
	return Date{
		year:  year,
		month: month,
		day:   day,
	}
}

// Format is YYYY-MM-DD
func NewDateFromString(str string) (Date, error) {
	parts := strings.Split(str, "-")

	year, yearErr := strconv.Atoi(parts[0])
	if yearErr != nil {
		return dateDef(), yearErr
	}
	month, monthErr := strconv.Atoi(parts[1])
	if monthErr != nil {
		return dateDef(), monthErr
	}
	day, dayErr := strconv.Atoi(parts[2])
	if dayErr != nil {
		return dateDef(), dayErr
	}

	return Date{
		year:  year,
		month: month,
		day:   day,
	}, nil
}

func (date Date) Year() int {
	return date.year
}

func (date Date) Month() int {
	return date.month
}

func (date Date) Day() int {
	return date.day
}

// ToZinc representes the object as: "YYYY-MM-DD"
func (date Date) ToZinc() string {
	return date.encode()
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
