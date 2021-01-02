package haystack

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Coord models a geographic coordinate as latitude and longitude
type Coord struct {
	lat float64
	lng float64
}

// NewCoord creates a new Coord object. lat/lng are clipped to reasonable values.
func NewCoord(lat float64, lng float64) *Coord {
	lat = math.Min(math.Max(lat, -90.0), 90.0) // Clip to -90 to 90
	lng = math.Min(math.Max(lng, 0.0), 180)    // Clip to 0 to 180
	return &Coord{lat: lat, lng: lng}
}

// Lat returns the latitude.
func (coord *Coord) Lat() float64 {
	return coord.lat
}

// Lng returns the longitude.
func (coord *Coord) Lng() float64 {
	return coord.lng
}

// ToZinc representes the object as: "C(<lat>,<lng>)"
func (coord *Coord) ToZinc() string {
	result := "C("
	result = result + fmt.Sprintf("%g", coord.lat)
	result = result + ","
	result = result + fmt.Sprintf("%g", coord.lng)
	result = result + ")"

	return result
}

// MarshalJSON representes the object as: "c:<lat>,<lng>"
func (coord *Coord) MarshalJSON() ([]byte, error) {
	result := "c:" + fmt.Sprintf("%g", coord.lat) + "," + fmt.Sprintf("%g", coord.lng)
	return json.Marshal(result)
}

// UnmarshalJSON interprets the json value: "c:<lat>,<lng>"
func (coord *Coord) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(jsonStr, "c:") {
		return errors.New("Input value does not begin with c:")
	}
	coordSplit := strings.Split(jsonStr[2:len(jsonStr)], ",")

	lat, latErr := strconv.ParseFloat(coordSplit[0], 64)
	if latErr != nil {
		return latErr
	}
	lng, lngErr := strconv.ParseFloat(coordSplit[1], 64)
	if lngErr != nil {
		return lngErr
	}

	*coord = *NewCoord(lat, lng)

	return nil
}

// MarshalHayson representes the object as: "{\"_kind\":\"coord\",\"lat\":<lat>,\"lng\":<lng>}"
func (coord *Coord) MarshalHayson() ([]byte, error) {
	builder := new(strings.Builder)
	builder.WriteString("{\"_kind\":\"coord\",\"lat\":")
	builder.WriteString(fmt.Sprintf("%g", coord.lat))
	builder.WriteString(",\"lng\":")
	builder.WriteString(fmt.Sprintf("%g", coord.lng))
	builder.WriteString("}")
	return []byte(builder.String()), nil
}
