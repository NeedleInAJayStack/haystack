package haystack

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"time"
)

// DateTime models a timestamp with a specific timezone.
type DateTime struct {
	date     *Date
	time     *Time
	tzOffset int    // offset in seconds from UTC
	tz       string // IANA database city name
}

// NewDateTime creates a new DateTime object. The values are not validated for correctness.
func NewDateTime(date *Date, time *Time, tzOffset int, tz string) *DateTime {
	return &DateTime{
		date:     date,
		time:     time,
		tzOffset: tzOffset,
		tz:       tz,
	}
}

// NewDateTimeRaw creates a new DateTime object. The values are not validated for correctness.
func NewDateTimeRaw(year int, month int, day int, hour int, min int, sec int, ms int, tzOffset int, tz string) *DateTime {
	date := NewDate(year, month, day)
	time := NewTime(hour, min, sec, ms)
	return &DateTime{
		date:     date,
		time:     time,
		tzOffset: tzOffset,
		tz:       tz,
	}
}

// NewDateTimeFromString creates a DateTime object from a string in the format: "YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func NewDateTimeFromString(str string) (*DateTime, error) {
	var input scanner.Scanner
	input.Init(strings.NewReader(str))
	curRune := input.Next()

	dateStr := strings.Builder{}
	for curRune != 'T' && curRune != scanner.EOF {
		dateStr.WriteRune(curRune)
		curRune = input.Next()
	}
	date, dateErr := NewDateFromIso(dateStr.String())
	if dateErr != nil {
		return dateTimeDef(), dateErr
	}

	curRune = input.Next() // Skip over 'T'

	timeStr := strings.Builder{}
	for curRune != '-' && curRune != '+' && curRune != 'Z' && curRune != scanner.EOF {
		timeStr.WriteRune(curRune)
		curRune = input.Next()
	}
	time, timeErr := NewTimeFromIso(timeStr.String())
	if timeErr != nil {
		return dateTimeDef(), timeErr
	}

	tz := "UTC"
	tzOffset := 0
	if curRune == '-' || curRune == '+' { // In this case we have an offset specified
		neg := curRune == '-'
		curRune = input.Next() // Skip over '+' or '-'

		hourStr := strings.Builder{}
		for curRune != ':' && curRune != scanner.EOF {
			hourStr.WriteRune(curRune)
			curRune = input.Next()
		}
		hour, hourErr := strconv.Atoi(hourStr.String())
		if hourErr != nil {
			return dateTimeDef(), hourErr
		}

		curRune = input.Next() // Skip over ':'

		minStr := strings.Builder{}
		for curRune != ' ' && curRune != scanner.EOF {
			minStr.WriteRune(curRune)
			curRune = input.Next()
		}
		min, minErr := strconv.Atoi(minStr.String())
		if minErr != nil {
			return dateTimeDef(), minErr
		}
		tzOffset = hour*3600 + min*60

		curRune = input.Next() // Skip over ' '

		tzStr := strings.Builder{}
		for curRune != scanner.EOF {
			tzStr.WriteRune(curRune)
			curRune = input.Next()
		}
		tz = tzStr.String()

		if neg {
			tzOffset = tzOffset * -1
		}
	}
	// Otherwise it's UTC

	return NewDateTime(date, time, tzOffset, tz), nil
}

func NewDateTimeFromGo(goTime time.Time) *DateTime {
	hDate := NewDate(
		goTime.Year(),
		int(goTime.Month()),
		goTime.Day(),
	)
	hTime := NewTime(
		goTime.Hour(),
		goTime.Minute(),
		goTime.Second(),
		goTime.Nanosecond()/1000,
	)
	location := goTime.Location()
	hTz := "UTC"
	if location != time.UTC {
		tzName := goTime.Location().String()
		hTz = strings.Split(tzName, "/")[1] // Don't include the region, only the city.
	}
	_, hTzOffset := goTime.Zone()

	return NewDateTime(hDate, hTime, hTzOffset, hTz)
}

func dateTimeDef() *DateTime {
	return NewDateTime(&Date{}, &Time{}, 0, "UTC")
}

// Date returns the date of the object.
func (dateTime *DateTime) Date() *Date {
	return dateTime.date
}

// Time returns the date of the object.
func (dateTime *DateTime) Time() *Time {
	return dateTime.time
}

// Tz returns the timezone of the object.
func (dateTime *DateTime) Tz() string {
	return dateTime.tz
}

// TzOffset returns the timezone offset of the object.
func (dateTime *DateTime) TzOffset() int {
	return dateTime.tzOffset
}

// ToZinc representes the object as: "YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func (dateTime *DateTime) ToZinc() string {
	buf := strings.Builder{}
	dateTime.encodeTo(&buf, true)
	return buf.String()
}

// MarshalJSON representes the object as: "t:YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func (dateTime *DateTime) MarshalJSON() ([]byte, error) {
	buf := strings.Builder{}
	buf.WriteString("t:")
	dateTime.encodeTo(&buf, true)
	return json.Marshal(buf.String())
}

// UnmarshalJSON interprets the json value: "t:YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func (dateTime *DateTime) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newDateTime, newErr := dateTimeFromJSON(jsonStr)
	*dateTime = *newDateTime
	return newErr
}

func dateTimeFromJSON(jsonStr string) (*DateTime, error) {
	if !strings.HasPrefix(jsonStr, "t:") {
		return nil, errors.New("Input value does not begin with 't:'")
	}
	dateTimeStr := jsonStr[2:]

	parseDateTime, parseErr := NewDateTimeFromString(dateTimeStr)
	if parseErr != nil {
		return nil, parseErr
	}
	return parseDateTime, nil
}

// MarshalHayson representes the object as: "{\"_kind\":\"dateTime\",\"val\":\"YYYY-MM-DD'T'hh:mm:ss.FFFz\",\"tz\":\"zzzz\"}"
func (dateTime *DateTime) MarshalHayson() ([]byte, error) {
	buf := strings.Builder{}
	buf.WriteString("{\"_kind\":\"dateTime\",\"val\":\"")
	dateTime.encodeTo(&buf, false)
	buf.WriteString("\",\"tz\":\"")
	buf.WriteString(dateTime.tz)
	buf.WriteString("\"}")
	return []byte(buf.String()), nil
}

func (dateTime *DateTime) encodeTo(buf *strings.Builder, includeTz bool) {
	buf.WriteString(dateTime.date.toIso())
	buf.WriteRune('T')
	buf.WriteString(dateTime.time.toIso())
	if dateTime.tzOffset == 0 {
		buf.WriteRune('Z')
	} else {
		offset := dateTime.tzOffset
		if offset < 0 {
			buf.WriteRune('-')
			offset = offset * -1
		} else {
			buf.WriteRune('+')
		}
		hr := offset / 3600
		min := (offset % 3600) / 60

		if hr < 10 {
			buf.WriteRune('0')
		}
		buf.WriteString(fmt.Sprintf("%d", hr))
		buf.WriteRune(':')

		if min < 10 {
			buf.WriteRune('0')
		}
		buf.WriteString(fmt.Sprintf("%d", min))
	}
	if includeTz {
		buf.WriteRune(' ')
		buf.WriteString(dateTime.tz)
	}
}
