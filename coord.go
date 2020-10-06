package haystack

import (
	"fmt"
	"math"
)

// Coord models a geographic coordinate as latitude and longitude
type Coord struct {
	lat float64
	lng float64
}

// NewCoord creates a new Coord object. lat/lng are clipped to reasonable values.
func NewCoord(lat float64, lng float64) Coord {
	lat = math.Min(math.Max(lat, -90.0), 90.0) // Clip to -90 to 90
	lng = math.Min(math.Max(lng, 0.0), 180)    // Clip to 0 to 180
	return Coord{lat: lat, lng: lng}
}

// Lat returns the latitude.
func (coord Coord) Lat() float64 {
	return coord.lat
}

// Lng returns the longitude.
func (coord Coord) Lng() float64 {
	return coord.lng
}

// ToZinc representes the object as: "C(<lat>,<lng>)"
func (coord Coord) ToZinc() string {
	result := "C("
	result = result + fmt.Sprintf("%g", coord.lat)
	result = result + ","
	result = result + fmt.Sprintf("%g", coord.lng)
	result = result + ")"

	return result
}
