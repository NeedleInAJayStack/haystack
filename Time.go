package haystack

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Time models a time of day tag value.
type Time struct {
	hour int
	min  int
	sec  int
	ms   int
}

// NewTime creates a new Time object. The values are not validated for correctness.
func NewTime(hour int, min int, sec int, ms int) Time {
	return Time{
		hour: hour,
		min:  min,
		sec:  sec,
		ms:   ms,
	}
}

// NewTimeFromString creates a Time object from a string in the format: "hh:mm:ss" or "hh:mm:ss.mmm"
func NewTimeFromString(str string) (Time, error) {
	parts := strings.Split(str, ":")

	hour, hourErr := strconv.Atoi(parts[0])
	if hourErr != nil {
		return Time{}, hourErr
	}
	min, minErr := strconv.Atoi(parts[1])
	if minErr != nil {
		return Time{}, minErr
	}

	sec := 0
	ms := 0
	if strings.Contains(parts[2], ".") { // Split ms out if included
		secParts := strings.Split(parts[2], ".")

		secVal, secErr := strconv.Atoi(secParts[0])
		if secErr != nil {
			return Time{}, secErr
		}
		sec = secVal
		msPart := secParts[1]
		msVal, msErr := strconv.Atoi(msPart)
		if msErr != nil {
			return Time{}, msErr
		}
		// Support inputting up to 3 digit accuracy
		if len(msPart) == 1 {
			ms = msVal * 100
		} else if len(msPart) == 2 {
			ms = msVal * 10
		} else if len(msPart) == 3 {
			ms = msVal
		} else {
			return Time{}, errors.New("ms section contained more than 3 digits")
		}
	} else {
		secVal, secErr := strconv.Atoi(parts[2])
		if secErr != nil {
			return Time{}, secErr
		}
		sec = secVal
	}

	return Time{
		hour: hour,
		min:  min,
		sec:  sec,
		ms:   ms,
	}, nil
}

// Hour returns the hours of the object.
func (time Time) Hour() int {
	return time.hour
}

// Min returns the minutes of the object.
func (time Time) Min() int {
	return time.min
}

// Sec returns the seconds of the object.
func (time Time) Sec() int {
	return time.sec
}

// Millis returns the milliseconds of the object.
func (time Time) Millis() int {
	return time.ms
}

// ToZinc representes the object as: "hh:mm:ss.FFF"
func (time Time) ToZinc() string {
	return time.encode()
}

// ToJSON representes the object as: "h:hh:mm:ss.FFF"
func (time Time) ToJSON() string {
	return "h:"+time.encode()
}

func (time Time) encode() string {
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

func (time1 Time) equals(time2 *Time) bool {
	return time1.hour == time2.hour &&
		time1.min == time2.min &&
		time1.sec == time2.sec &&
		time1.ms == time2.ms
}
