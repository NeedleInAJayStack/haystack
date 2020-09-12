package haystack

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Time struct {
	hour int
	min  int
	sec  int
	ms   int // Optional
}

func timeDef() Time {
	return Time{
		hour: 0,
		min:  0,
		sec:  0,
		ms:   0,
	}
}

// Format is hh:mm:ss or hh:mm:ss.mmm
func timeFromStr(str string) (Time, error) {
	parts := strings.Split(str, ":")

	hour, hourErr := strconv.Atoi(parts[0])
	if hourErr != nil {
		return timeDef(), hourErr
	}
	min, minErr := strconv.Atoi(parts[1])
	if minErr != nil {
		return timeDef(), minErr
	}

	sec := 0
	ms := 0
	if strings.Contains(parts[2], ".") { // Split ms out if included
		secParts := strings.Split(parts[2], ".")

		secVal, secErr := strconv.Atoi(secParts[0])
		if secErr != nil {
			return timeDef(), secErr
		}
		sec = secVal
		msPart := secParts[1]
		msVal, msErr := strconv.Atoi(msPart)
		if msErr != nil {
			return timeDef(), msErr
		}
		// Support inputting up to 3 digit accuracy
		if len(msPart) == 1 {
			ms = msVal * 100
		} else if len(msPart) == 2 {
			ms = msVal * 10
		} else if len(msPart) == 3 {
			ms = msVal
		} else {
			return timeDef(), errors.New("ms section contained more than 3 digits")
		}
	} else {
		secVal, secErr := strconv.Atoi(parts[2])
		if secErr != nil {
			return timeDef(), secErr
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

// Format is hh:mm:ss.mmm
func (time Time) toZinc() string {
	return time.encode()
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
