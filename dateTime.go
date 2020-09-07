package haystack

import (
	"fmt"
	"strings"
	"time"
)

type DateTime struct {
	date     Date
	time     Time
	tz       string // IANA database city name
	tzOffset int    // offset in seconds from UTC
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

// Format is YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz
func (dateTime *DateTime) toZinc() string {
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
