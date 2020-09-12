package haystack

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"time"
)

type DateTime struct {
	date     Date
	time     Time
	tz       string // IANA database city name
	tzOffset int    // offset in seconds from UTC
}

// Parses date from "YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func dateTimeFromStr(str string) (DateTime, error) {
	var input scanner.Scanner
	input.Init(strings.NewReader(str))
	curRune := input.Next()

	dateStr := strings.Builder{}
	for curRune != 'T' && curRune != scanner.EOF {
		dateStr.WriteRune(curRune)
		curRune = input.Next()
	}
	date, dateErr := dateFromStr(dateStr.String())
	if dateErr != nil {
		return dateTimeDef(), dateErr
	}

	curRune = input.Next() // Skip over 'T'

	timeStr := strings.Builder{}
	for curRune != '-' && curRune != '+' && curRune != 'Z' && curRune != scanner.EOF {
		timeStr.WriteRune(curRune)
		curRune = input.Next()
	}
	time, timeErr := timeFromStr(timeStr.String())
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

	return DateTime{
		date:     date,
		time:     time,
		tz:       tz,
		tzOffset: tzOffset,
	}, nil
}

func dateTimeFromGo(goTime time.Time) DateTime {
	hDate := Date{
		year:  goTime.Year(),
		month: int(goTime.Month()),
		day:   goTime.Day(),
	}
	hTime := Time{
		hour: goTime.Hour(),
		min:  goTime.Minute(),
		sec:  goTime.Second(),
		ms:   goTime.Nanosecond() / 1000,
	}
	location := goTime.Location()
	hTz := "UTC"
	if location != time.UTC {
		tzName := goTime.Location().String()
		hTz = strings.Split(tzName, "/")[1] // Don't include the region, only the city.
	}
	_, hTzOffset := goTime.Zone()

	return DateTime{
		date:     hDate,
		time:     hTime,
		tz:       hTz,
		tzOffset: hTzOffset,
	}
}

func dateTimeDef() DateTime {
	return DateTime{
		date:     dateDef(),
		time:     timeDef(),
		tz:       "UTC",
		tzOffset: 0,
	}
}

// Format is YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz
func (dateTime DateTime) toZinc() string {
	buf := strings.Builder{}
	dateTime.encodeTo(&buf)
	return buf.String()
}

func (dateTime *DateTime) encodeTo(buf *strings.Builder) {
	buf.WriteString(dateTime.date.encode())
	buf.WriteRune('T')
	buf.WriteString(dateTime.time.encode())
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
	buf.WriteRune(' ')
	buf.WriteString(dateTime.tz)
}
