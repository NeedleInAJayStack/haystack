package haystack

import (
	"encoding/json"
	"errors"
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
func NewDate(year int, month int, day int) *Date {
	return &Date{
		year:  year,
		month: month,
		day:   day,
	}
}

// NewDateFromIso creates a Date object from a string in the format: "YYYY-MM-DD"
func NewDateFromIso(str string) (*Date, error) {
	parts := strings.Split(str, "-")

	year, yearErr := strconv.Atoi(parts[0])
	if yearErr != nil {
		return NewDate(0, 0, 0), yearErr
	}
	month, monthErr := strconv.Atoi(parts[1])
	if monthErr != nil {
		return NewDate(0, 0, 0), monthErr
	}
	day, dayErr := strconv.Atoi(parts[2])
	if dayErr != nil {
		return NewDate(0, 0, 0), dayErr
	}

	return NewDate(year, month, day), nil
}

// Year returns the years of the object.
func (date *Date) Year() int {
	return date.year
}

// Month returns the numerical month of the object.
func (date *Date) Month() int {
	return date.month
}

// Day returns the day-of-month of the object.
func (date *Date) Day() int {
	return date.day
}

// ToZinc representes the object as: "YYYY-MM-DD"
func (date *Date) ToZinc() string {
	return date.toIso()
}

// MarshalJSON representes the object as: "d:YYYY-MM-DD"
func (date *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal("d:" + date.toIso())
}

// UnmarshalJSON interprets the json value: "d:YYYY-MM-DD"
func (date *Date) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(jsonStr, "d:") {
		return errors.New("Input value does not begin with d:")
	}
	dateStr := jsonStr[2:len(jsonStr)]

	parseDate, parseErr := NewDateFromIso(dateStr)
	if parseErr != nil {
		return parseErr
	}
	*date = *parseDate

	return nil
}

// MarshalHayson representes the object as: "{\"_kind\":\"date\",\"val\":\"YYYY-MM-DD\""}"
func (date *Date) MarshalHayson() ([]byte, error) {
	return []byte("{\"_kind\":\"date\",\"val\":\"" + date.toIso() + "\"}"), nil
}

func (date *Date) toIso() string {
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
