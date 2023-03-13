package haystack

import (
	"encoding/json"
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

// NewTimeFromIso creates a Time object from a string in the format: "hh:mm:ss[.mmm]"
func NewTimeFromIso(str string) (Time, error) {
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
		msMaxIndex := len(msPart)
		if msMaxIndex > 3 {
			msMaxIndex = 3
		}
		msPart = msPart[:msMaxIndex]
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
		}
		// It shouldn't be higher because we clamp it
	} else {
		secVal, secErr := strconv.Atoi(parts[2])
		if secErr != nil {
			return Time{}, secErr
		}
		sec = secVal
	}

	return NewTime(hour, min, sec, ms), nil
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

// ToZinc representes the object as: "hh:mm:ss[.mmm]"
func (time Time) ToZinc() string {
	return time.toIso()
}

// MarshalJSON representes the object as: "h:hh:mm:ss[.mmm]"
func (time Time) MarshalJSON() ([]byte, error) {
	return json.Marshal("h:" + time.toIso())
}

// UnmarshalJSON interprets the json value: "h:hh:mm:ss[.mmm]"
func (time *Time) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newTime, newErr := timeFromJSON(jsonStr)
	*time = newTime
	return newErr
}

func timeFromJSON(jsonStr string) (Time, error) {
	if !strings.HasPrefix(jsonStr, "h:") {
		return Time{}, errors.New("value does not begin with 'h:'")
	}
	timeStr := jsonStr[2:]

	return NewTimeFromIso(timeStr)
}

// MarshalHayson representes the object as: "{\"_kind\":\"time\",\"val\":\"hh:mm:ss[.mmm]\""}"
func (time Time) MarshalHayson() ([]byte, error) {
	return []byte("{\"_kind\":\"time\",\"val\":\"" + time.toIso() + "\"}"), nil
}

func (time Time) toIso() string {
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

func (time Time) equals(otherTime Time) bool {
	return time.hour == otherTime.hour &&
		time.min == otherTime.min &&
		time.sec == otherTime.sec &&
		time.ms == otherTime.ms
}
